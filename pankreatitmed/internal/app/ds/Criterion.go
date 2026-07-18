package ds

func (Criterion) TableName() string { return "criteria" }

type Criterion struct {
	ID          uint    `gorm:"primaryKey"`
	Code        string  `gorm:"type:varchar(12);unique;not null"` // №1..№11
	Name        string  `gorm:"type:varchar(120);not null"`
	Description string  `gorm:"type:text;not null"`
	Duration    string  `gorm:"type:varchar(64);not null"` // "1 календарный день"
	HomeVisit   bool    `gorm:"not null;default:false"`
	ImageURL    *string `gorm:"type:varchar(255)"`
	Status      string  `gorm:"type:varchar(12);not null;default:'active';check:status IN ('active','deleted')"`
	Unit        string  `gorm:"type:varchar(32);not null;default:''"` // единиуа измерения
	RefLow      *float64
	RefHigh     *float64
}
