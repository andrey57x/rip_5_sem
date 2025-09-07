package ds

type ReactionCalculation struct {
	ID int `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	ReactionID    int `gorm:"not null;uniqueIndex:idx_reaction_calculation"`
	CalculationID int `gorm:"not null;uniqueIndex:idx_reaction_calculation"`

	Amount int `gorm:"default:1"`

	Reaction    Reaction    `gorm:"foreignKey:ReactionID"`
	Calculation Calculation `gorm:"foreignKey:CalculationID"`
}
