package ds

import "time"

const (
	OrderStatusDraft     = "draft"
	OrderStatusDeleted   = "deleted"
	OrderStatusFormed    = "formed"
	OrderStatusCompleted = "completed"
	OrderStatusRejected  = "rejected"
)

func (PankreatitOrder) TableName() string { return "pankreatitorders" }

type PankreatitOrder struct {
	ID            uint      `gorm:"primaryKey"`
	Status        string    `gorm:"type:varchar(12);not null;check:status IN ('draft','deleted','formed','completed','rejected')"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	CreatorID     uint      `gorm:"not null"`
	Creator       MedUser   `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT;"`
	FormedAt      *time.Time
	FinishedAt    *time.Time
	ModeratorID   *uint
	Moderator     *MedUser `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT;"`
	RansonScore   *int
	MortalityRisk *string `gorm:"type:varchar(24)"`
}
