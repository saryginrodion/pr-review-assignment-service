package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type PullRequestResponse struct {
	PullRequest PullRequest `json:"pr"`
	ReplacedBy  *string      `json:"replaced_by,omitempty"`
}

func ToPullRequestResponse(pr entities.PullRequest) PullRequestResponse {
	return PullRequestResponse{
		PullRequest: ToPullRequest(pr),
		ReplacedBy:  nil,
	}
}

func ToPullRequestReassignedResponse(pr entities.PullRequest, oldReviewerID string) PullRequestResponse {
	return PullRequestResponse{
		PullRequest: ToPullRequest(pr),
		ReplacedBy:  &oldReviewerID,
	}
}
