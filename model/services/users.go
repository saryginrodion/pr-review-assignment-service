package services

import (
	"context"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsersService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewUsersService(db *gorm.DB, ctx context.Context) UsersService {
	return UsersService{
		db:  db,
		ctx: ctx,
	}
}

func (s *UsersService) SetIsActive(userId string, isActive bool) (*entities.User, error) {
	var user entities.User
	res := s.db.Model(&entities.User{}).
		Where("id = ?", userId).
		Update("is_active", isActive).
		Clauses(clause.Returning{}).
		Scan(&user)

	if res.RowsAffected == 0 {
		return nil, &ErrNotFound{}
	} else if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
