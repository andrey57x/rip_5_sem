package apitypes

import "Backend/internal/app/ds"

// ReactionCalculationJSON модель расчёта для реакции
// @ID ReactionCalculationJSON
type ReactionCalculationJSON struct {
	ID            int     `json:"id"`
	ReactionID    int     `json:"reaction_id"`
	CalculationID int     `json:"calculation_id"`
	OutputMass    float32 `json:"output_mass"`
	InputMass     float32 `json:"input_mass"`
}

func ReactionCalculationToJSON(reactionCalculation ds.ReactionCalculation) ReactionCalculationJSON {
	return ReactionCalculationJSON{
		ID:            reactionCalculation.ID,
		ReactionID:    reactionCalculation.ReactionID,
		CalculationID: reactionCalculation.CalculationID,
		OutputMass:    reactionCalculation.OutputMass,
		InputMass:     reactionCalculation.InputMass,
	}
}

func ReactionCalculationFromJSON(reactionCalculationJSON ReactionCalculationJSON) ds.ReactionCalculation {
	return ds.ReactionCalculation{
		OutputMass: reactionCalculationJSON.OutputMass,
	}
}
