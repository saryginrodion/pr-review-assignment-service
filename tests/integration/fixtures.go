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

	return db
}

func CleanUpDb(db *gorm.DB) {
	db.Migrator().DropTable(entities.User{}, entities.Team{})
}

// Setups team with name "TeamA" and user with id "user1"
func SetupTeamAUser1(db *gorm.DB, t *testing.T) {
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	members := []entities.User{
		{ID: "user1", Username: "Alice", IsActive: true},
	}

	team, err := teams.Create("TeamA", members)
	assert.NoError(t, err)
	assert.Equal(t, "TeamA", team.TeamName)
	assert.Len(t, team.Members, 1)
}

