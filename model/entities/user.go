package entities

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"not null"`
	TeamName string
	Team     Team `gorm:"foreignKey:TeamName;references:TeamName"`
	IsActive bool
}
