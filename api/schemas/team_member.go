package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required,min=1"`
	Username string `json:"username" validate:"required,min=1"`
	IsActive bool   `json:"is_active" validate:"required"`
}

func ToTeamMember(u entities.User) TeamMember {
	return TeamMember{
		UserID:   u.ID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}

func TeamMemberToUserModel(member TeamMember, teamName string) entities.User {
	return entities.User{
		ID:                   member.UserID,
		Username:             member.Username,
		TeamName:             teamName,
		IsActive:             member.IsActive,
	}
}

