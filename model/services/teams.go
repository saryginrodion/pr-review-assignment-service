package services

import (
	"context"
	"errors"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"gorm.io/gorm"
)

type TeamsService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewTeamsService(db *gorm.DB, ctx context.Context) TeamsService {
	return TeamsService{
		db:  db,
		ctx: ctx,
	}
}

func (s *TeamsService) Create(teamName string, members []entities.User) (*entities.Team, error) {
	tx := s.db.Begin()

	newTeam := entities.Team{
		Name: teamName,
	}

	// omitting members so we can create them manually.
	// if gorm handles it automatically, it will not only create members
	// but also update existing ones with matching IDs, which is not the intended behavior (we need to error if user already exists).
	err := gorm.
		G[entities.Team](tx).
		Omit("Members").
		Create(s.ctx, &newTeam)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		tx.Rollback()
		return nil, &ErrTeamExists{TeamName: teamName}
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	for i, member := range members {
		member.TeamName = teamName
		members[i] = member
	}

	if err := tx.Create(members).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		tx.Rollback()
		return nil, &ErrUserExistsInAnotherTeam{}
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	newTeam.Members = members

	return &newTeam, nil
}

func (s *TeamsService) Get(teamName string) (*entities.Team, error) {
	team, err := gorm.
		G[entities.Team](s.db).
		Where("name = ?", teamName).
		Preload("Members", nil).
		First(s.ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &ErrNotFound{}
	} else if err != nil {
		return nil, err
	}

	return &team, nil
}
