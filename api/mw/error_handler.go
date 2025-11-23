package mw

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/schemas"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/saryginrodion/stackable"
)

type ErrorMapping struct {
	ErrType reflect.Type
	Code    string
	Status  int
}

// Маппинг для известных ошибок (все они приходят с сервисов)
var ApiErrorMappings = []ErrorMapping{
	{reflect.TypeOf(&services.ErrUserExistsInAnotherTeam{}), "USER_EXISTS_IN_ANOTHER_TEAM", http.StatusConflict},
	{reflect.TypeOf(&services.ErrNotFound{}), "NOT_FOUND", http.StatusNotFound},
	{reflect.TypeOf(&services.ErrTeamExists{}), "TEAM_EXISTS", http.StatusConflict},
	{reflect.TypeOf(&services.ErrPullRequestExists{}), "PR_EXISTS", http.StatusConflict},
	{reflect.TypeOf(&services.ErrNoCandidates{}), "NO_CANDIDATE", http.StatusConflict},
	{reflect.TypeOf(&services.ErrPullRequestMerged{}), "PR_MERGED", http.StatusConflict},
}

func MapError(mappings []ErrorMapping, err error) *ErrorMapping {
	for _, e := range mappings {
		target := reflect.New(e.ErrType).Interface()
		if errors.As(err, target) {
			return &e
		}
	}
	return nil
}

var ErrorsHandlerMiddleware = stackable.WrapFunc(
	func(ctx *context.Context, next func() error) error {
		err := next()
		if err == nil {
			return nil
		}

		// Parsing error
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			ctx.Response, _ = stackable.JsonResponse(
				http.StatusBadRequest,
				schemas.NewErrorResponse("PARSE_FAIL", err.Error()),
			)
			return nil
		}

		// Validation error
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			ctx.Response, _ = stackable.JsonResponse(
				http.StatusUnprocessableEntity,
				schemas.NewErrorResponse("VALIDATION_FAIL", err.Error()),
			)
			return nil
		}

		// Проверка ошибок из маппинга
		knownErrorMapping := MapError(ApiErrorMappings, err)
		if knownErrorMapping != nil {
			ctx.Response, _ = stackable.JsonResponse(
				knownErrorMapping.Status,
				schemas.NewErrorResponse(knownErrorMapping.Code, err.Error()),
			)
			return nil
		}

		// unknown error
		ctx.Shared.Logger.Error("Unhandled error", "err", err.Error())
		ctx.Response, _ = stackable.JsonResponse(
			http.StatusInternalServerError,
			schemas.NewErrorResponse("UNKNOWN", err.Error()),
		)
		return err
	},
)
