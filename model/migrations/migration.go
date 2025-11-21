package migrations

import (
	"context"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, ctx context.Context) error {
	err := db.AutoMigrate(&entities.Team{}, &entities.User{})

	if err != nil {
		return err
	}

	return nil
}
