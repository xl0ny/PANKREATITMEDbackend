package ds

func (MedUser) TableName() string { return "medusers" }

// TODO создать роли
type MedUser struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"size:32;uniqueIndex;not null"`
	Password    string `gorm:"size:128;not null"`
	IsModerator bool   `gorm:"not null;default:false"`
}
