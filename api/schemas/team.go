package schemas

import (
	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
)

type Team struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,dive"`
}

func ToTeam(team entities.Team) Team {
	return Team{
		TeamName: team.Name,
		Members:  utils.MapSlice(ToTeamMember, team.Members),
	}
}
