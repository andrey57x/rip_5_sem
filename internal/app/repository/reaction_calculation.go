package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
)

func (r *Repository) GetReactionCalculations(calculationID int) ([]ds.ReactionCalculation, error) {
	var reactionCalculations []ds.ReactionCalculation
	err := r.db.Where("calculation_id = ?", calculationID).Find(&reactionCalculations).Error
	if err != nil {
		return nil, err
	}
	return reactionCalculations, nil
}

func (r *Repository) GetReactionCalculation(reactionID int, calculationID int) (ds.ReactionCalculation, error) {
	var reactionCalculation ds.ReactionCalculation
	err := r.db.Where("reaction_id = ? and calculation_id = ?", reactionID, calculationID).First(&reactionCalculation).Error
	if err != nil {
		return ds.ReactionCalculation{}, err
	}
	return reactionCalculation, nil
}

func (r *Repository) DeleteReactionFromCalculation(calculationID int, reactionID int) (ds.Calculation, error) {
	var calculation ds.Calculation
	err := r.db.Where("id = ?", calculationID).First(&calculation).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	err = r.db.Where("reaction_id = ? and calculation_id = ?", reactionID, calculationID).Delete(&ds.ReactionCalculation{}).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	return calculation, nil
}

func (r *Repository) ChangeReactionCalculation(calculationID int, reactionID int, reactionCalculationJSON apitypes.ReactionCalculationJSON) (ds.ReactionCalculation, error) {
	var reactionCalculation ds.ReactionCalculation
	err := r.db.Model(&reactionCalculation).Where("reaction_id = ? and calculation_id = ?", reactionID, calculationID).Updates(apitypes.ReactionCalculationFromJSON(reactionCalculationJSON)).First(&reactionCalculation).Error
	if err != nil {
		return ds.ReactionCalculation{}, err
	}
	return reactionCalculation, nil
}
