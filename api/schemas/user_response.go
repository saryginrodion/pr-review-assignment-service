package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type UserResponse struct {
	User User `json:"user"`
}

func ToUserResponse(user entities.User) UserResponse {
	return UserResponse{
		User: ToUser(user),
	}
}
