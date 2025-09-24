package repository

import (
	"Backend/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var ErrorNoDraft = errors.New("no draft for this user")

func (r *Repository) GetMassCalculationReactions(id int) ([]ds.ReactionInfo, ds.MassCalculation, error) {
	var calculation ds.MassCalculation
	err := r.db.Where("id = ?", id).First(&calculation).Error
	if err != nil {
		return []ds.ReactionInfo{}, ds.MassCalculation{}, err
	} else if calculation.Status == "deleted" {
		return []ds.ReactionInfo{}, ds.MassCalculation{}, errors.New("you can`t watch deleted calculations")
	}

	var reactions []ds.Reaction
	var reactionCalculations []ds.ReactionCalculation
	sub := r.db.Table("reaction_calculations").Where("calculation_id = ?", calculation.ID).Find(&reactionCalculations)
	err = r.db.Where("id IN (?)", sub.Select("reaction_id")).Find(&reactions).Error

	var reactionInfos []ds.ReactionInfo
	for _, reaction := range reactions {
		for _, reactionCalculation := range reactionCalculations {
			if reaction.ID == reactionCalculation.ReactionID {
				reactionInfos = append(reactionInfos, ds.ReactionInfo{
					ID:                 reaction.ID,
					Title:              reaction.Title,
					Reagent:            reaction.Reagent,
					Product:            reaction.Product,
					ConversationFactor: reaction.ConversationFactor,
					ImgLink:            reaction.ImgLink,

					OutputMass: reactionCalculation.OutputMass,
					InputMass:  reactionCalculation.InputMass,
				})
				break
			}
		}
	}

	if err != nil {
		return []ds.ReactionInfo{}, ds.MassCalculation{}, err
	}

	return reactionInfos, calculation, nil
}

// GetCartCount для получения количества услуг в заявке
func (r *Repository) GetCartCount() int {
	var count int64
	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	calculation, err := r.CheckCurrentMassCalculationDraft(creatorID)
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.ReactionCalculation{}).Where("calculation_id = ?", calculation.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting reactions in reaction_calculations:", err)
	}

	return int(count)
}

func (r *Repository) CheckCurrentMassCalculationDraft(creatorID int) (ds.MassCalculation, error) {
	var calculation ds.MassCalculation

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&calculation)
	if res.Error != nil {
		return ds.MassCalculation{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.MassCalculation{}, ErrorNoDraft
	}
	return calculation, nil
}

func (r *Repository) GetMassCalculationDraft(creatorID int) (ds.MassCalculation, error) {
	calculation, err := r.CheckCurrentMassCalculationDraft(creatorID)
	if err == ErrorNoDraft {
		calculation = ds.MassCalculation{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&calculation)
		if result.Error != nil {
			return ds.MassCalculation{}, result.Error
		}
		return calculation, nil
	} else if err != nil {
		return ds.MassCalculation{}, err
	}
	return calculation, nil
}

func (r *Repository) DeleteMassCalculation(id int) error {
	return r.db.Exec("UPDATE calculations SET status = 'deleted' WHERE id = ?", id).Error
}
