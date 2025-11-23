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

	t.Log("PR:", pr)
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

func TestPullRequestMerge(t *testing.T) {
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

	// Проверяем что реквест открылся
	pr, err := prs.Create("pr-merge", "Merge Test", *author)
	assert.NoError(t, err)
	assert.Equal(t, entities.PULL_REQUEST_OPEN, pr.Status)

	// Проверяем что реквест смерджился
	mergedPR, err := prs.Merge("pr-merge")
	assert.NoError(t, err)
	assert.Equal(t, entities.PULL_REQUEST_MERGED, mergedPR.Status)
	assert.True(t, mergedPR.MergedAt.Valid)
	// Проверяем время, чтоб ПР смерджился +- в пределах 2 секунд
	assert.WithinDuration(t, time.Now(), mergedPR.MergedAt.Time, time.Second*2)

	// Проверяем что реквест все это действительно теперь лежит в БД
	pr, err = prs.GetFull("pr-merge")
	assert.NoError(t, err)
	assert.Equal(t, entities.PULL_REQUEST_MERGED, pr.Status)
	assert.True(t, pr.MergedAt.Valid)

	// Проверяем что у assignedreviewers изменилось время LastAssignedAt
	assert.Equal(t, len(pr.AssignedReviewers), 2)
	assert.True(t, pr.AssignedReviewers[0].LastAssignedAt.Valid)
	assert.WithinDuration(t, pr.AssignedReviewers[0].LastAssignedAt.Time, time.Now(), time.Second*4)
}

func TestPullRequestMergeNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	prs := services.NewPullRequestsService(db, context.Background())

	_, err := prs.Merge("does_not_exist")
	assert.Error(t, err)
}
