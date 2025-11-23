package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"gorm.io/gorm"
)

type PullRequestsService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewPullRequestsService(db *gorm.DB, ctx context.Context) PullRequestsService {
	return PullRequestsService{
		db:  db,
		ctx: ctx,
	}
}

func (s *PullRequestsService) GetFull(pullRequestID string) (*entities.PullRequest, error) {
	var pr entities.PullRequest

	err := s.db.
		Preload("Author").
		Preload("AssignedReviewers").
		First(&pr, "id = ?", pullRequestID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &ErrNotFound{}
	} else if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *PullRequestsService) GetUsersToAssign(
	teamName string,
	reviewersCount int,
	excludeUserIDs []string,
) ([]entities.User, error) {
	var users []entities.User

	s.db.
		Model(&entities.User{}).
		Where("team_name = ?", teamName).
		Where("is_active = ?", true).
		Where("id NOT IN ?", excludeUserIDs).
		Order("last_assigned_at ASC NULLS FIRST").
		Limit(reviewersCount).
		Find(&users)

	if s.db.Error != nil {
		return nil, s.db.Error
	}

	if len(users) != reviewersCount {
		return nil, &ErrNoCandidates{}
	}

	return users, nil
}

func (s *PullRequestsService) Create(pullRequestID string, pullRequestName string, author entities.User) (*entities.PullRequest, error) {
	tx := s.db.Begin()

	newPR := entities.PullRequest{
		ID:       pullRequestID,
		Name:     pullRequestName,
		AuthorID: author.ID,
		Author:   author,
		Status:   entities.PULL_REQUEST_OPEN,
	}

	// Создаем PR без AssignedReviewers и без изменения Author
	err := tx.Omit("Author").Create(&newPR).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		tx.Rollback()
		return nil, &ErrPullRequestExists{
			PullRequestID: pullRequestID,
		}
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Выбираем ревьюверов
	prServiceTx := NewPullRequestsService(tx, s.ctx)
	usersToAssign, err := prServiceTx.GetUsersToAssign(
		newPR.Author.TeamName,
		2,
		[]string{author.ID},
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Им присваивем time.Now() как last_assigned_at, и сохраняем их IDs для апдейта
	now := time.Now()
	usersToAssignIDs := make([]string, len(usersToAssign))
	for i, user := range usersToAssign {
		usersToAssign[i].LastAssignedAt.Time = now
		usersToAssignIDs[i] = user.ID
	}

	// Апдейтим
	err = tx.
		Model(&entities.User{}).
		Where("id in ?", usersToAssignIDs).
		Update("last_assigned_at", now).
		Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Присваеваем PR'у ревьюверов
	err = tx.
		Model(&newPR).
		Association("AssignedReviewers").
		Append(usersToAssign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &newPR, err
}

func (s *PullRequestsService) Merge(pullRequestID string) (*entities.PullRequest, error) {
	err := s.db.Model(&entities.PullRequest{}).
		Where("id = ?", pullRequestID).
		Updates(&entities.PullRequest{
			MergedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Status: entities.PULL_REQUEST_MERGED,
		},
		).Error

	if err != nil {
		return nil, err
	}

	pr, err := s.GetFull(pullRequestID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
