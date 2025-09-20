package repository

import (
	"errors"
	"fmt"

	"Backend/internal/app/ds"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("is_delete = false").Find(&reactions).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("title ILIKE ? and is_delete = ?", "%"+title+"%", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) AddReactionToCalculation(calculationID int, reactionID int) error {
	var reaction ds.Reaction
	if err := r.db.First(&reaction, reactionID).Error; err != nil {
		return err
	}

	var calculation ds.Calculation
	if err := r.db.First(&calculation, calculationID).Error; err != nil {
		return err
	}
	reactionCalculation := ds.ReactionCalculation{}
	result := r.db.Where("reaction_id = ? and calculation_id = ?", reactionID, calculationID).Find(&reactionCalculation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.ReactionCalculation{
		ReactionID:    reactionID,
		CalculationID: calculationID,
	}).Error
}

func CalculateMass(mass float32, conversationFactor, outputKoef float32) (float32, error) {
	if conversationFactor == 0 {
		return 0, errors.New("invalid conversation factor")
	}
	if outputKoef == 0 || outputKoef > 1 {
		return 0, errors.New("invalid output koeficient")
	}
	return mass / conversationFactor / outputKoef, nil
}
