package integration_test

import (
	"context"
	"testing"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/model/migrations"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open("user=tests password=tests host=localhost port=15432 dbname=tests sslmode=disable"), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = migrations.Migrate(db, context.Background())
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	db = db.Debug()

	return db
}

func CleanUpDB(db *gorm.DB) {
	db.Migrator().DropTable(
		entities.User{},
		entities.Team{},
		entities.PullRequest{},
		"pr_assigned_reviewers",
	)
}

func SetupTeamAndUsers(
	db *gorm.DB,
	t *testing.T,
	teamName string,
	users []entities.User,
) {
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	team, err := teams.Create(teamName, users)
	assert.NoError(t, err)
	assert.Equal(t, teamName, team.Name)
	assert.Len(t, team.Members, len(users))
}

