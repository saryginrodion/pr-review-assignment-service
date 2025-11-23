package schemas

import (
	"time"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
)

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id" validate:"required"`
	PullRequestName   string     `json:"pull_request_name" validate:"required"`
	AuthorID          string     `json:"author_id" validate:"required"`
	Status            string     `json:"status" validate:"required,oneof=OPEN MERGED"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func ToPullRequest(pr entities.PullRequest) PullRequest {
	var mergedAt *time.Time
	if pr.MergedAt.Valid {
		mergedAt = &pr.MergedAt.Time
	}

	return PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: utils.MapSlice(func(user entities.User) string { return user.ID }, pr.AssignedReviewers),
		CreatedAt:         pr.CreatedAt,
		MergedAt:          mergedAt,
	}
}
