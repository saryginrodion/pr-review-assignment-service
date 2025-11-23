package schemas

import "github.com/saryginrodion/pr_review_assignment_service/model/entities"

type TeamResponse struct {
	Team Team `json:"team"`
}

func ToTeamResponse(team entities.Team) TeamResponse {
	return TeamResponse{
		Team: ToTeam(team),
	}
}
