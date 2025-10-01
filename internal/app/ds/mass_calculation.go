package ds

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MassCalculation struct {
	ID          int           `gorm:"primaryKey;autoIncrement"`
	OutputKoef  float32       `gorm:"float;default:1"`
	Status      string        `gorm:"type:varchar(16);not null"`
	DateCreate  time.Time     `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID   uuid.UUID     `gorm:"not null"`
	ModeratorID uuid.NullUUID `gorm:"default:null"`

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
