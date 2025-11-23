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
	assert.Equal(t, len(mergedPR.AssignedReviewers), 2)
	assert.True(t, mergedPR.AssignedReviewers[0].LastAssignedAt.Valid)
	assert.WithinDuration(t, mergedPR.AssignedReviewers[0].LastAssignedAt.Time, time.Now(), time.Second*4)

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

func TestPullRequestReassignSuccess(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "author", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
		{ID: "rev3", Username: "rev3", TeamName: "TeamA", IsActive: true},
		{ID: "rev4", Username: "rev3", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	// Создаем PR
	pr, err := prs.Create("pr-reassign", "PR Reassign", *author)
	assert.NoError(t, err)
	assert.Len(t, pr.AssignedReviewers, 2)

	oldReviewerIDs := []string{pr.AssignedReviewers[0].ID}

	// Делаем Reassign
	updatedPR, err := prs.Reassign("pr-reassign", oldReviewerIDs)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(updatedPR.AssignedReviewers))

	// Проверяем, что старый ревьювер заменен новым
	foundOld := false
	for _, r := range updatedPR.AssignedReviewers {
		if r.ID == oldReviewerIDs[0] {
			foundOld = true
		}
	}
	assert.False(t, foundOld, "Old reviewer should be replaced")

	for _, r := range updatedPR.AssignedReviewers {
		assert.NotEqual(t, "author", r.ID)
		assert.True(t, r.LastAssignedAt.Valid)
		assert.WithinDuration(t, time.Now(), r.LastAssignedAt.Time, time.Second*2)
	}
}

func TestPullRequestReassignNoCandidates(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "author", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	// Создаем PR с одним ревьювером (его потом попытаемся переназначить)
	pr, err := prs.Create("pr-no-candidates", "PR No Candidates", *author)
	assert.NoError(t, err)

	oldReviewerIDs := []string{pr.AssignedReviewers[0].ID}

	// Попытка Reassign должна вернуть ErrNoCandidates
	_, err = prs.Reassign("pr-no-candidates", oldReviewerIDs)
	assert.Error(t, err)
	_, ok := err.(*services.ErrNoCandidates)
	assert.True(t, ok)
}

func TestPullRequestReassignNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	prs := services.NewPullRequestsService(db, context.Background())

	_, err := prs.Reassign("non-existent-pr", []string{"rev1"})
	assert.Error(t, err)
	_, ok := err.(*services.ErrNotFound)
	assert.True(t, ok)
}

func TestPullRequestReassignMerged(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanUpDB(db)

	SetupTeamAndUsers(db, t, "TeamA", []entities.User{
		{ID: "author", Username: "author", TeamName: "TeamA", IsActive: true},
		{ID: "rev1", Username: "rev1", TeamName: "TeamA", IsActive: true},
		{ID: "rev2", Username: "rev2", TeamName: "TeamA", IsActive: true},
		{ID: "rev3", Username: "rev3", TeamName: "TeamA", IsActive: true},
		{ID: "rev4", Username: "rev3", TeamName: "TeamA", IsActive: true},
	})

	users := services.NewUsersService(db, context.Background())
	author, err := users.Get("author")
	assert.NoError(t, err)

	prs := services.NewPullRequestsService(db, context.Background())

	// Создаем PR
	pr, err := prs.Create("pr-reassign", "PR Reassign", *author)
	assert.NoError(t, err)
	assert.Len(t, pr.AssignedReviewers, 2)

	oldReviewerIDs := []string{pr.AssignedReviewers[0].ID}

	// Мерджим
	pr, err = prs.Merge(pr.ID)
	assert.NoError(t, err)

	// Пытаемся сделать reassign
	_, err = prs.Reassign("pr-reassign", oldReviewerIDs)
	assert.Error(t, err)

	_, ok := err.(*services.ErrPullRequestMerged)
	assert.True(t, ok)
}
