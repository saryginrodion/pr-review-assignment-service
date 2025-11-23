package services

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
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

	usersToAssign, err = prServiceTx.Assign(pullRequestID, usersToAssign)
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
	pr, err := s.GetFull(pullRequestID)
	if err != nil {
		return nil, err
	}

	if pr.Status == entities.PULL_REQUEST_MERGED {
		return pr, nil
	}

	pr.MergedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	pr.Status = entities.PULL_REQUEST_MERGED

	err = s.db.Save(pr).Error

	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PullRequestsService) Assign(pullRequestID string, usersToAssign []entities.User) ([]entities.User, error) {
	err := s.db.Model(&entities.PullRequest{ID: pullRequestID}).
		Association("AssignedReviewers").
		Append(usersToAssign)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	usersToAssignIDs := make([]string, len(usersToAssign))
	for i, user := range usersToAssign {
		usersToAssign[i].LastAssignedAt.Time = now
		usersToAssignIDs[i] = user.ID
	}

	// Апдейтим
	err = s.db.
		Model(&entities.User{}).
		Where("id in ?", usersToAssignIDs).
		Update("last_assigned_at", now).
		Error
	if err != nil {
		return nil, err
	}

	return usersToAssign, nil
}

func (s *PullRequestsService) Reassign(pullRequestID string, oldReviewerIDs []string) (*entities.PullRequest, error) {
	// Получаем полный PR
	pr, err := s.GetFull(pullRequestID)
	if err != nil {
		return nil, err
	}

	// Pull request не смерджен
	if pr.Status == entities.PULL_REQUEST_MERGED {
		return nil, &ErrPullRequestMerged{}
	}


	// Проверяем, что ревьювер действительно есть в assigned reviewers
	assignedReviwersIDs := utils.MapSlice(
		func(rev entities.User) string { return rev.ID },
		pr.AssignedReviewers,
	)

	filteredOldReviewerIDs := utils.FilterSlice(
		func(id string) bool { return slices.Contains(assignedReviwersIDs, id) },
		oldReviewerIDs,
	)

	if len(filteredOldReviewerIDs) != len(oldReviewerIDs) {
		return nil, &ErrNotFound{}
	}

	tx := s.db.Begin()

	// Удаляем старых ревьюверов
	var oldReviewers []entities.User
	for _, userID := range oldReviewerIDs {
		oldReviewers = append(oldReviewers, entities.User{ID: userID})
	}
	err = tx.Model(&pr).Association("AssignedReviewers").Delete(oldReviewers)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Составляем ignore список для реассайна
	excludeUserIDs := make([]string, len(oldReviewerIDs))
	// Игнорируем старых ревьюверов
	copy(excludeUserIDs, oldReviewerIDs)
	// Игнорируем автора
	excludeUserIDs = append(excludeUserIDs, pr.AuthorID)
	// Игнорируем других уже назначенных пользователей
	for _, user := range pr.AssignedReviewers {
		if !slices.Contains(oldReviewerIDs, user.ID) {
			excludeUserIDs = append(excludeUserIDs, user.ID)
		}
	}

	// Получаем список реассигнации
	prServiceTx := NewPullRequestsService(tx, s.ctx)
	usersToAssign, err := prServiceTx.GetUsersToAssign(
		pr.Author.TeamName,
		len(oldReviewerIDs),
		excludeUserIDs,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	usersToAssign, err = prServiceTx.Assign(pullRequestID, usersToAssign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Получаем PR с обновленными ассоциациями
	pr, err = s.GetFull(pullRequestID)
	if err != nil {
		return nil, err
	}

	return pr, tx.Error
}
