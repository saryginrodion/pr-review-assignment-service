package api

import (
	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/mw"
	"github.com/saryginrodion/stackable"
)

func NewStack(sharedState *context.SharedState) stackable.Stackable[context.SharedState, context.LocalState] {
	stack := stackable.NewStackable[context.SharedState, context.LocalState](sharedState)
	stack.AddHandler(&mw.ErrorsHandlerMiddleware)
	stack.AddHandler(&mw.LoggingMiddleware)
	return stack
}

