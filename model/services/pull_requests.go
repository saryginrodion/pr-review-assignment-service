package services

import (
	"context"
	"errors"

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

func (s *PullRequestsService) GetUsersToAssign(
	teamName string,
	reviewersCount int,
	excludeUserIDs []string,
) (*[]entities.User, error) {
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

	return &users, nil
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

	err := tx.Model(&entities.PullRequest{}).Create(&newPR).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		tx.Rollback()
		return nil, &ErrPullRequestExists{
			PullRequestID: pullRequestID,
		}
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}


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

	newPR.AssignedReviewers = *usersToAssign
	err = tx.Save(&newPR).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &newPR, nil
}
