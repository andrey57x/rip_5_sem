package ds

import (
	"database/sql"
	"time"
)

type Calculation struct {
	ID          int           `gorm:"primaryKey;autoIncrement"`
	OutputKoef  float32       `gorm:"float;default:1"`
	Status      string        `gorm:"type:varchar(16);not null"`
	DateCreate  time.Time     `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID   int           `gorm:"not null"`
	ModeratorID sql.NullInt64 `gorm:"default:null"`

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
