package routes

import (
	"net/http"

	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/stackable"
)

var GetIndex = stackable.WrapFunc(
	func(ctx *context.Context, next func() error) error {
		if ctx.Request.URL.Path != "/" || ctx.Request.Method != "GET" {
			ctx.Response = stackable.NewHttpResponse(
				http.StatusNotFound,
				"text/html",
				"<h1>404 - Not found</h1>",
			)

			return next()
		}

		ctx.Response = stackable.NewHttpResponse(
			http.StatusOK,
			"text/html",
			"<h1>Hello from Backend!</h1>",
		)

		return next()
	},
)
