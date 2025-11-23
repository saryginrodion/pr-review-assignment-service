package schemas

type PullRequestMergeBody struct {
	PullRequetsID string `json:"pull_request_id" validate:"required,min=1"`
}
