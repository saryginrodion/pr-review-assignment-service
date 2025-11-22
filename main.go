package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/saryginrodion/pr_review_assignment_service/env"
	"github.com/saryginrodion/pr_review_assignment_service/model/migrations"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelInfo}))
	env := env.Env()

	db, err := gorm.Open(postgres.Open(env.POSTGRES_DSN), &gorm.Config{
		TranslateError: true,
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

	teams := services.NewTeamsService(db, ctx)
	team, err := teams.Get("test_team0")
	// team, err := teams.Create("test_team0", []entities.User{
	// 	{
	// 		ID:       "u1",
	// 		Username: "Alice",
	// 		IsActive: true,
	// 	},
	// 	{
	// 		ID:       "u2",
	// 		Username: "Bob",
	// 		IsActive: true,
	// 	},
	// 	{
	// 		ID:       "u3",
	// 		Username: "Gorilka",
	// 		IsActive: false,
	// 	},
	// })

	if err != nil {
		logger.Error("HUI", "err", err.Error())
	} else {
		logger.Info("Team:", "team", team)
	}
}
