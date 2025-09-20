package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var ErrorNotAllowed = errors.New("you are not allowed")
var ErrorNoDraft = errors.New("no draft for this user")

func (r *Repository) GetCalculationReactions(id int) ([]ds.Reaction, ds.Calculation, error) {
	calculation, err := r.GetSingleCalculation(id)
	if err != nil {
		return []ds.Reaction{}, ds.Calculation{}, err
	}

	var reactions []ds.Reaction
	sub := r.db.Table("reaction_calculations").Where("calculation_id = ?", calculation.ID)
	err = r.db.Where("id IN (?)", sub.Select("reaction_id")).Find(&reactions).Error

	if err != nil {
		return []ds.Reaction{}, ds.Calculation{}, err
	}

	return reactions, calculation, nil
}

// GetCartCount для получения количества услуг в заявке
func (r *Repository) GetCartCount(creatorID int) int {
	var count int64

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
		return ds.Calculation{}, ErrorNoDraft
	}
	return calculation, nil
}

func (r *Repository) GetCalculationDraft(creatorID int) (ds.Calculation, bool, error) {
	calculation, err := r.CheckCurrentCalculationDraft(creatorID)
	if err == ErrorNoDraft {
		calculation = ds.Calculation{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&calculation)
		if result.Error != nil {
			return ds.Calculation{}, false, result.Error
		}
		return calculation, true, nil
	} else if err != nil {
		return ds.Calculation{}, false, err
	}
	return calculation, false, nil
}

func (r *Repository) DeleteCalculation(id int) error {
	return r.db.Exec("UPDATE calculations SET status = 'deleted' WHERE id = ?", id).Error
}

func (r *Repository) GetCalculations(from, to time.Time, status string) ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("date_create > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_create < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Find(&calculations).Error
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

func (r *Repository) ChangeCalculation(id int, calculationJSON apitypes.CalculationJSON) (ds.Calculation, error) {
	calculation := ds.Calculation{}
	if id < 0 {
		return ds.Calculation{}, errors.New("invalid id, it must be >= 0")
	}
	if calculationJSON.OutputKoef <= 0 || calculationJSON.OutputKoef > 1 {
		return ds.Calculation{}, errors.New("invalid output koeficient")
	}
	err := r.db.Where("id = ? and status != 'deleted'", id).First(&calculation).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	err = r.db.Model(&calculation).Updates(apitypes.CalculationFromJSON(calculationJSON)).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	return calculation, nil
}

func (r *Repository) GetSingleCalculation(id int) (ds.Calculation, error) {
	if id < 0 {
		return ds.Calculation{}, errors.New("invalid id, it must be >= 0")
	}
	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Calculation{}, err
	}
	var calculation ds.Calculation
	err = r.db.Where("id = ?", id).First(&calculation).Error
	if err != nil {
		return ds.Calculation{}, err
	// } else if user.ID != calculation.CreatorID && !user.IsModerator {
	// 	return ds.Calculation{}, ErrorNotAllowed
	} else if calculation.Status == "deleted" && !user.IsModerator {
		return ds.Calculation{}, errors.New("calculation is deleted")
	}
	return calculation, nil
}

func (r *Repository) FormCalculation(id int, status string) (ds.Calculation, error) {
	calculation, err := r.GetSingleCalculation(id)
	if err != nil {
		return ds.Calculation{}, err
	}

	if calculation.Status != "draft" {
		return ds.Calculation{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&calculation).Updates(ds.Calculation{
		Status: status,
		DateForm: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Calculation{}, err
	}

	return calculation, nil
}

func (r *Repository) ModerateCalculation(id int, status string) (ds.Calculation, error) {
	if status != "completed" && status != "rejected" {
		return ds.Calculation{}, errors.New("wrong status")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Calculation{}, err
	}

	if !user.IsModerator {
		return ds.Calculation{}, errors.New("you are not a moderator")
	}

	calculation, err := r.GetSingleCalculation(id)
	if err != nil {
		return ds.Calculation{}, err
	} else if calculation.Status != "formed" {
		return ds.Calculation{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&calculation).Updates(ds.Calculation{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Calculation{}, err
	}

	if status == "completed" {
		reactionCalculations, err := r.GetReactionCalculations(calculation.ID)
		if err != nil {
			return ds.Calculation{}, err
		}
		for _, reactionCalculation := range reactionCalculations {
			reaction, err := r.GetReaction(reactionCalculation.ReactionID)
			if err != nil {
				return ds.Calculation{}, err
			}
			mass, err := CalculateMass(reactionCalculation.OutputMass, reaction.ConversationFactor, calculation.OutputKoef)
			if err != nil {
				return ds.Calculation{}, err
			}
			err = r.db.Model(&reactionCalculation).Updates(ds.ReactionCalculation{
				InputMass: mass,
			}).Error
			if err != nil {
				return ds.Calculation{}, err
			}
		}
	}

	return calculation, nil
}
