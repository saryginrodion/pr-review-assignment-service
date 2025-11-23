package context

import (
	"log/slog"

	"github.com/saryginrodion/stackable"
	"github.com/saryginrodion/stackable/middleware"
	"gorm.io/gorm"
)

type SharedState struct {
	DB *gorm.DB
	Logger *slog.Logger
}

type LocalState struct {
	*middleware.LocalRequestId
}

func (l LocalState) Default() any {
	return LocalState{
		LocalRequestId: &middleware.LocalRequestId{},
	}
}

type Context = stackable.Context[SharedState, LocalState]
