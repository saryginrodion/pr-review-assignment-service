package schemas

type PullRequestReassignBody struct {
	PullRequestID string `json:"pull_request_id" validate:"required,min=1"`
	OldReviewerID string `json:"old_reviewer_id" validate:"required,min=1"`
}
