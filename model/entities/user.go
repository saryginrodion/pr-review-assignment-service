package entities

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	// username 
	IsActive bool      `gorm:""`
}
