package migrations

import (
	"context"

	"github.com/saryginrodion/pr_review_assignment_service/model/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, ctx context.Context) error {
	db.DisableForeignKeyConstraintWhenMigrating = true
	err := db.AutoMigrate(&entities.Team{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		return err
	}

	if !db.Migrator().HasConstraint(&entities.User{}, "TeamName") {
		db.Migrator().CreateConstraint(&entities.User{}, "TeamName")
	}

	return nil
}
