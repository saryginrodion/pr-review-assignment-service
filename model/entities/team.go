package entities

type Team struct {
	TeamName string `gorm:"primaryKey;not null"`
	Members  []User `gorm:"foreignKey:TeamName;references:TeamName;constraint:OnDelete:CASCADE"`
}
