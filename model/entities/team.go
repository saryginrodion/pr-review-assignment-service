package entities

type Team struct {
	Name string `gorm:"primaryKey;not null"`
	Members  []User `gorm:"foreignKey:TeamName;references:Name;constraint:OnDelete:CASCADE"`
}
