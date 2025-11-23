package routes

import (
	"net/http"

	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/schemas"
	apiUtils "github.com/saryginrodion/pr_review_assignment_service/api/utils"
	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/saryginrodion/pr_review_assignment_service/utils"
	"github.com/saryginrodion/stackable"
)

var TeamAdd = stackable.WrapFunc(
	func (ctx *context.Context, next func() error) error {
		body, err := apiUtils.ParseAndValidateJson(ctx.Request.Body, schemas.Team{})
		if err != nil {
			return err
		}

		teams := services.NewTeamsService(ctx.Shared.DB, ctx.Request.Context())
		members := utils.MapSlice(
			func (teamMember schemas.TeamMember) entities.User { return schemas.TeamMemberToUserModel(teamMember, body.TeamName) },
			body.Members,
		)
		team, err := teams.Create(body.TeamName, members)

		if err != nil {
			return err
		}

		ctx.Response, _ = stackable.JsonResponse(
			http.StatusCreated,
			schemas.ToTeam(*team),
		)

		return next()
	},
)

var TeamGet = stackable.WrapFunc(
	func (ctx *context.Context, next func() error) error {
		teamName := ctx.Request.URL.Query().Get("team_name")
		teams := services.NewTeamsService(ctx.Shared.DB, ctx.Request.Context())

		team, err := teams.Get(teamName)
		if err != nil {
			return err
		}

		ctx.Response, _ = stackable.JsonResponse(
			http.StatusOK,
			schemas.ToTeam(*team),
		)
		
		return next()
	},
)
