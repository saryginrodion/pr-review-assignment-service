package routes

import (
	"net/http"

	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/schemas"
	apiUtils "github.com/saryginrodion/pr_review_assignment_service/api/utils"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/saryginrodion/stackable"
)

var PullRequestCreate = stackable.WrapFunc(
	func(ctx *context.Context, next func() error) error {
		body, err := apiUtils.ParseAndValidateJson(ctx.Request.Body, schemas.PullRequestCreateBody{})
		if err != nil {
			return err
		}

		users := services.NewUsersService(ctx.Shared.DB, ctx.Request.Context())
		author, err := users.Get(body.AuthorID)
		if err != nil {
			return err
		}

		prs := services.NewPullRequestsService(ctx.Shared.DB, ctx.Request.Context())
		pr, err := prs.Create(body.PullRequestID, body.PullRequestName, *author)
		if err != nil {
			return err
		}

		ctx.Response, _ = stackable.JsonResponse(
			http.StatusOK,
			schemas.ToPullRequest(*pr),
		)

		return next()
	},
)

var PullRequestMerge = stackable.WrapFunc(
	func(ctx *context.Context, next func() error) error {
		body, err := apiUtils.ParseAndValidateJson(ctx.Request.Body, schemas.PullRequestMergeBody{})
		if err != nil {
			return err
		}

		prs := services.NewPullRequestsService(ctx.Shared.DB, ctx.Request.Context())
		pr, err := prs.Merge(body.PullRequetsID)
		if err != nil {
			return err
		}

		ctx.Response, _ = stackable.JsonResponse(
			http.StatusOK,
			schemas.ToPullRequest(*pr),
		)

		return next()
	},
)
