package integration_test

import (
	"context"
	"testing"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/stretchr/testify/assert"
)


func TestTeamsService_CreateAndGet(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	members := []entities.User{
		{ID: "user1", Username: "Alice", IsActive: true},
		{ID: "user2", Username: "Bob", IsActive: true},
	}

	team, err := teams.Create("TeamA", members)
	assert.NoError(t, err)
	assert.Equal(t, "TeamA", team.Name)
	assert.Len(t, team.Members, 2)

	team, err = teams.Get("TeamA")
	assert.NoError(t, err)
	assert.Equal(t, "TeamA", team.Name)
	assert.Len(t, team.Members, 2)
}

func TestTeamsService_CreateTeamAlreadyExists(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	members := []entities.User{
		{ID: "user1", Username: "Alice", IsActive: true},
	}

	_, err := teams.Create("TeamA", members)
	assert.NoError(t, err)

	// Try creating the same team again
	_, err = teams.Create("TeamA", members)
	assert.Error(t, err)
	_, ok := err.(*services.ErrTeamExists)
	assert.True(t, ok)
}

func TestTeamsService_CreateUserAlreadyExistsInAnotherTeam(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	_, err := teams.Create("TeamA", []entities.User{
		{ID: "user1", Username: "Alice", IsActive: true},
	})
	assert.NoError(t, err)

	_, err = teams.Create("TeamB", []entities.User{
		{ID: "user1", Username: "Alice", IsActive: true},
	})
	assert.Error(t, err)
	_, ok := err.(*services.ErrUserExistsInAnotherTeam)
	assert.True(t, ok)
}

func TestTeamsService_GetNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)
	ctx := context.Background()
	teams := services.NewTeamsService(db, ctx)

	_, err := teams.Get("NonExistentTeam")
	assert.Error(t, err)
	_, ok := err.(*services.ErrNotFound)
	assert.True(t, ok)
}
