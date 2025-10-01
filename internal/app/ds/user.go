package ds

import "github.com/google/uuid"

type User struct {
	UUID        uuid.UUID `gorm:"primary_key"`
	Login       string    `gorm:"type:varchar(25);unique;not null"`
	Password    string    `gorm:"type:varchar(100);not null"`
	IsModerator bool      `gorm:"type:boolean;default:false"`
}
