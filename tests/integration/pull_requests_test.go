package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"github.com/saryginrodion/pr_review_assignment_service/model/services"
	"github.com/stretchr/testify/assert"
)

func TestPullRequestCreate(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "author", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
		{ID: "rev3", Username: "rev3", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())


	pr, err := prs.Create("pr1", "PR Name", *author)
	assert.NoError(t, err)
	assert.Equal(t, "pr1", pr.ID)
	assert.Equal(t, "PR Name", pr.Name)
	assert.Equal(t, entities.PULL_REQUEST_OPEN, pr.Status)

	assert.Len(t, pr.AssignedReviewers, 2)

	for _, r := range pr.AssignedReviewers {
		assert.NotEqual(t, "author", r.ID)
	}
}

func TestPullRequestCreateDuplicate(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "author", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
		{ID: "rev3", Username: "rev3", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	_, err = prs.Create("pr1", "PR1", *author)
	assert.NoError(t, err)

	_, err = prs.Create("pr1", "PR1", *author)
	assert.Error(t, err)

	_, ok := err.(*services.ErrPullRequestExists)
	assert.True(t, ok)
}

func TestPullRequestCreateNoCandidates(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "a", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	_, err = prs.Create("pr1", "name", *author)
	assert.Error(t, err)

	_, ok := err.(*services.ErrNoCandidates)
	assert.True(t, ok)
}

func TestPullRequestBeforeSaveConstraintMerged(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "a", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	pr, err := prs.Create("pr1", "name", *author)
	assert.NoError(t, err)

	pr.Status = entities.PULL_REQUEST_MERGED
	err = db.Save(pr).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MergedAt cannot be null")

	pr.MergedAt.Valid = true
	pr.MergedAt.Time = time.Now()

	err = db.Save(pr).Error
	assert.NoError(t, err)
}

func TestPullRequestCreateNoCandidatesInactive(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "a", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: false},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: false},
		{ID: "rev3", Username: "rev3", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	_, err = prs.Create("pr1", "name", *author)
	assert.Error(t, err)

	_, ok := err.(*services.ErrNoCandidates)
	assert.True(t, ok)
}
