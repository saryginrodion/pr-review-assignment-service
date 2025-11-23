package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/saryginrodion/pr_review_assignment_service/api"
	"github.com/saryginrodion/pr_review_assignment_service/api/context"
	"github.com/saryginrodion/pr_review_assignment_service/api/swaggerui"
	"github.com/saryginrodion/pr_review_assignment_service/env"
	"github.com/saryginrodion/pr_review_assignment_service/model/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelInfo}))
	env := env.Env()

	// Setting up DB
	db, err := gorm.Open(postgres.Open(env.POSTGRES_DSN), &gorm.Config{
		TranslateError: true,
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
			slogGorm.WithTraceAll(),
			slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
		),
	})
	if err != nil {
		slog.Error("Error on DB connection: ", "err", err)
		os.Exit(1)
	}

	if err := migrations.Migrate(db, db.Statement.Context); err != nil {
		slog.Error("Failed on migrations: ", "err", err)
		os.Exit(1)
	}

	stack := api.NewStack(&context.SharedState{
		DB: db,
		Logger: logger,
	})
	httpServer := api.HttpServer(stack, ":" + strconv.Itoa(env.APP_PORT))
	swaggerui.SetupSwaggerUI()
	logger.Info("Starting server on :8000")
	logger.Error("Error on ListenAndServe", "err", httpServer.ListenAndServe())
}
