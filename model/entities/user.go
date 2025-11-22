package entities

type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	TeamName string `gorm:"not null;<-:create"`
	Team     Team   `gorm:"foreignKey:TeamName;references:Name"`
	IsActive bool   `gorm:"not null"`
}
