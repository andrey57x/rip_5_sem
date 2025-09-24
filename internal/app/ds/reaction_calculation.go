package ds

type ReactionCalculation struct {
	ID int `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	ReactionID    int `gorm:"not null;uniqueIndex:idx_reaction_calculation"`
	CalculationID int `gorm:"not null;uniqueIndex:idx_reaction_calculation"`

	OutputMass float32 `gorm:"float;not null;default:1000"`
	InputMass  float32 `gorm:"float;default:null"`

	Reaction    Reaction        `gorm:"foreignKey:ReactionID"`
	Calculation MassCalculation `gorm:"foreignKey:CalculationID"`
}
