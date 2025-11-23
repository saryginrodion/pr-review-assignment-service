package schemas

type PullRequestCreateBody struct {
	PullRequestID   string `json:"pull_request_id" validate:"required,min=1"`
	PullRequestName string `json:"pull_request_name" validate:"required,min=1"`
	AuthorID        string `json:"author_id" validate:"required,min=1"`
}
