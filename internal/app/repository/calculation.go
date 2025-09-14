package repository

import (
	"Backend/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var noDraftError = errors.New("no draft for this user")

func (r *Repository) GetCalculationReactions(id int) ([]ds.ReactionInfo, ds.Calculation, error) {

	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	var calculation ds.Calculation
	err := r.db.Where("id = ?", id).First(&calculation).Error
	if err != nil {
		return []ds.ReactionInfo{}, ds.Calculation{}, err
	} else if creatorID != calculation.CreatorID {
		return []ds.ReactionInfo{}, ds.Calculation{}, errors.New("you are not allowed")
	} else if calculation.Status == "deleted" {
		return []ds.ReactionInfo{}, ds.Calculation{}, errors.New("you can`t watch deleted calculations")
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
		return []ds.ReactionInfo{}, ds.Calculation{}, err
	}

	return reactionInfos, calculation, nil
}

// GetCartCount для получения количества услуг в заявке
func (r *Repository) GetCartCount() int {
	var count int64
	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	calculation, err := r.CheckCurrentCalculationDraft(creatorID)
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.ReactionCalculation{}).Where("calculation_id = ?", calculation.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting reactions in reaction_calculations:", err)
	}

	return int(count)
}

func (r *Repository) CheckCurrentCalculationDraft(creatorID int) (ds.Calculation, error) {
	var calculation ds.Calculation

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&calculation)
	if res.Error != nil {
		return ds.Calculation{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Calculation{}, noDraftError
	}
	return calculation, nil
}

func (r *Repository) GetCalculationDraft(creatorID int) (ds.Calculation, error) {
	calculation, err := r.CheckCurrentCalculationDraft(creatorID)
	if err == noDraftError {
		calculation = ds.Calculation{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&calculation)
		if result.Error != nil {
			return ds.Calculation{}, result.Error
		}
		return calculation, nil
	} else if err != nil {
		return ds.Calculation{}, err
	}
	return calculation, nil
}

func (r *Repository) DeleteCalculation(id int) error {
	return r.db.Exec("UPDATE calculations SET status = 'deleted' WHERE id = ?", id).Error
}
