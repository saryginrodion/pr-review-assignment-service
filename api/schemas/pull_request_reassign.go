package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type PullRequestReassignBody struct {
	PullRequestID string `json:"pull_request_id" validate:"required,min=1"`
	OldReviewerID string `json:"old_reviewer_id" validate:"required,min=1"`
}

type PullRequestReassignResponse struct {
	PullRequest PullRequest `json:"pr"`
	ReplacedBy  string      `json:"replaced_by"`
}

func ToPullRequestReassignResponse(pr entities.PullRequest, oldReviewerID string) PullRequestReassignResponse {
	return PullRequestReassignResponse{
		PullRequest: ToPullRequest(pr),
		ReplacedBy:  oldReviewerID,
	}
}
