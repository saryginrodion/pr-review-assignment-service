package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
	Status          string `json:"status" validate:"required,oneof=OPEN MERGED"`
}

func ToPullRequestShort(pr entities.PullRequest) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}
