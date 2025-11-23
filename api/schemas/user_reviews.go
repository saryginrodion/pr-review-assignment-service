package schemas

import (
	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
)

type UserReviews struct {
	UserID string `json:"string_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

func ToUserReviews(u entities.User) UserReviews {
	return UserReviews{
		UserID:   u.ID,
		PullRequests: utils.MapSlice(ToPullRequestShort, u.AssignedPullRequests),
	}
}
