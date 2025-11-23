package integration_test

import (
	"context"
	"testing"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/stretchr/testify/assert"
)

func TestUserGet(t *testing.T) {
	db := SetupTestDB(t)
	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ ID: "user1", Username: "username", IsActive: true },
	})
	defer CleanUpDB(db)

	users := services.NewUsersService(db, context.Background())
	user, err := users.Get("user1")
	assert.NoError(t, err)
	assert.Equal(t, "user1", user.ID)
}

func TestUserNotFound(t *testing.T) {
	db := SetupTestDB(t)
	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ ID: "user1", Username: "username", IsActive: true },
	})
	defer CleanUpDB(db)

	users := services.NewUsersService(db, context.Background())
	_, err := users.Get("user_not_exists")
	assert.Error(t, err)

	_, ok := err.(*services.ErrNotFound)
	assert.True(t, ok)
}

func TestUserSetActive(t *testing.T) {
	db := SetupTestDB(t)
	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ ID: "user1", Username: "username", IsActive: true },
	})
	defer CleanUpDB(db)


	users := services.NewUsersService(db, context.Background())

	user, err := users.SetIsActive("user1", false)
	assert.NoError(t, err)
	assert.Equal(t, false, user.IsActive)
	assert.Equal(t, "user1", user.ID)

	user, err = users.SetIsActive("user1", true)
	assert.NoError(t, err)
	assert.Equal(t, true, user.IsActive)
	assert.Equal(t, "user1", user.ID)
}

func TestUserSetActiveNotExists(t *testing.T) {
	db := SetupTestDB(t)
	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ ID: "user1", Username: "username", IsActive: true },
	})
	defer CleanUpDB(db)

	users := services.NewUsersService(db, context.Background())

	_, err := users.SetIsActive("user_that_does_not_exists", false)
	assert.Error(t, err)
	_, ok := err.(*services.ErrNotFound)
	assert.True(t, ok)
}
