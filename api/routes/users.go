package routes

import (
	"net/http"

	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/schemas"
	apiUtils "github.com/saryginrodion/pr_review_assignment_service/api/utils"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/saryginrodion/stackable"
)

var UserSetIsActive = stackable.WrapFunc(
	func(ctx *context.Context, next func() error) error {
		body, err := apiUtils.ParseAndValidateJson(ctx.Request.Body, schemas.UserSetIsActiveBody{})
		if err != nil {
			return err
		}

		users := services.NewUsersService(ctx.Shared.DB, ctx.Request.Context())
		user, err := users.SetIsActive(body.UserID, body.IsActive)
		if err != nil {
			return err
		}

		ctx.Response, _ = stackable.JsonResponse(
			http.StatusOK,
			schemas.ToUser(*user),
		)

		return next()
	},
)
