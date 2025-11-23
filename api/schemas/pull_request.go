package schemas

import (
	"database/sql"
	"time"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
)

type PullRequest struct {
	PullRequestID     string       `json:"pull_request_id" validate:"required"`
	PullRequestName   string       `json:"pull_request_name" validate:"required"`
	AuthorID          string       `json:"author_id" validate:"required"`
	Status            string       `json:"status" validate:"required,oneof=OPEN MERGED"`
	AssignedReviewers []string     `json:"assigned_reviewers"`
	CreatedAt         *time.Time   `json:"createdAt,omitempty"`
	MergedAt          sql.NullTime `json:"mergedAt"`
}

func ToPullRequest(pr entities.PullRequest) PullRequest {
	var createdAt *time.Time
	if !pr.CreatedAt.IsZero() {
		t := pr.CreatedAt
		createdAt = &t
	}

	return PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: utils.MapSlice(func(user entities.User) string { return user.ID }, pr.AssignedReviewers),
		CreatedAt:         createdAt,
		MergedAt:          pr.MergedAt,
	}
}
