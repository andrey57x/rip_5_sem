package ds

import (
	"database/sql"
	"time"
)

type Calculation struct {
	ID            int       `gorm:"primaryKey"`
	OutputMass    float32   `gorm:"float;"`
	OutputPercent float32   `gorm:"float;"`
	Status        string      `gorm:"type:varchar(16);not null"`
	DateCreate    time.Time `gorm:"not null"`
	DateUpdate    time.Time
	DateFinish    sql.NullTime `gorm:"default:null"`
	CreatorID     int          `gorm:"not null"`
	ModeratorID   int

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
