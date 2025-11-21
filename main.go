package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/saryginrodion/pr_review_assignment_service/env"
	"github.com/saryginrodion/pr_review_assignment_service/model/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, nil))
	env := env.Env()

	db, err := gorm.Open(postgres.Open(env.POSTGRES_DSN), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
			slogGorm.WithTraceAll(),
			slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
		),
	})

	if err != nil {
		logger.Error("Error on DB connection", "err", err)
	}

	ctx := context.Background()
	migrations.Migrate(db, ctx)
}
