package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type User struct {
	UserID   string `json:"user_id" validate:"required,min=1"`
	Username string `json:"username" validate:"required,min=1"`
	TeamName string `json:"team_name" validate:"required,min=1"`
	IsActive bool   `json:"is_active"`
}

func ToUser(u entities.User) User {
	return User{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
