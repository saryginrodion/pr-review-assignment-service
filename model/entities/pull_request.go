package entities

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PullRequest struct {
	ID                string `gorm:"primaryKey"`
	Name              string `gorm:"not null"`
	AuthorID          string `gorm:"not null;<-:create"`
	Author            User
	Status            PullRequestStatus `gorm:"not null;default:OPEN"`
	AssignedReviewers []User            `gorm:"many2many:pr_assigned_reviewers"`
	CreatedAt         time.Time
	MergedAt          sql.NullTime
}

func (pr *PullRequest) BeforeSave(tx *gorm.DB) (err error) {
	if pr.Status == PULL_REQUEST_MERGED && !pr.MergedAt.Valid {
		return errors.New("MergedAt cannot be null when status is MERGED")
	}

	return nil
}
