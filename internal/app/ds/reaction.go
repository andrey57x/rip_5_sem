package ds

type Reaction struct {
	ID                 int     `gorm:"primaryKey"`
	Title              string  `gorm:"type: varchar(64); not null"`
	Formula            string  `gorm:"type: varchar(32)"`
	ConversationFactor float32 `gorm:"type: float; not null"`
	ImgLink            string  `gorm:"type: varchar(32)"`
	Description        string  `gorm:"type: varchar(512)"`
	IsDelete           bool    `gorm:"type:boolean;default:false;not null"`
}
