package mw

import (
	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/stackable"
)

var LoggingMiddleware = stackable.WrapFunc(
	func(context *context.Context, next func() error) error {
		err := next()

		context.
			Shared.
			Logger.
			Info("Processed Request",
				"id", context.Local.RequestId(),
				"status", context.Response.Status(),
				"method", context.Request.Method,
				"path", context.Request.URL.Path,
				"err", err)

		return err
	},
)
