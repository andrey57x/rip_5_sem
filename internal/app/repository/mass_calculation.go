package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrorNotAllowed = errors.New("you are not allowed")
	ErrorNoDraft    = errors.New("no draft for this user")
	ErrorNotFound   = errors.New("not found")
)

func (r *Repository) GetMassCalculationReactions(id int) ([]ds.Reaction, ds.MassCalculation, error) {
	calculation, err := r.GetSingleMassCalculation(id)
	if err != nil {
		return []ds.Reaction{}, ds.MassCalculation{}, err
	}

	var reactions []ds.Reaction
	sub := r.db.Table("reaction_calculations").Where("calculation_id = ?", calculation.ID)
	err = r.db.Where("id IN (?)", sub.Select("reaction_id")).Find(&reactions).Error

	if err != nil {
		return []ds.Reaction{}, ds.MassCalculation{}, ErrorNotFound
	}

	return reactions, calculation, nil
}

// GetCartCount для получения количества услуг в заявке
func (r *Repository) GetCartCount(creatorID int) int {
	var count int64

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

func (r *Repository) GetMassCalculationDraft(creatorID int) (ds.MassCalculation, bool, error) {
	calculation, err := r.CheckCurrentMassCalculationDraft(creatorID)
	if err == ErrorNoDraft {
		calculation = ds.MassCalculation{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&calculation)
		if result.Error != nil {
			return ds.MassCalculation{}, false, result.Error
		}
		return calculation, true, nil
	} else if err != nil {
		return ds.MassCalculation{}, false, err
	}
	return calculation, false, nil
}

func (r *Repository) GetMassCalculations(from, to time.Time, status string) ([]ds.MassCalculation, error) {
	var calculations []ds.MassCalculation
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

func (r *Repository) ChangeMassCalculation(id int, calculationJSON apitypes.MassCalculationJSON) (ds.MassCalculation, error) {
	calculation := ds.MassCalculation{}
	if id < 0 {
		return ds.MassCalculation{}, errors.New("invalid id, it must be >= 0")
	}
	if *calculationJSON.OutputKoef <= 0 || *calculationJSON.OutputKoef > 1 {
		return ds.MassCalculation{}, errors.New("invalid output koeficient")
	}
	err := r.db.Where("id = ? and status != 'deleted'", id).First(&calculation).Error
	if err != nil {
		return ds.MassCalculation{}, ErrorNotFound
	}
	err = r.db.Model(&calculation).Updates(apitypes.MassCalculationFromJSON(calculationJSON)).Error
	if err != nil {
		return ds.MassCalculation{}, err
	}
	return calculation, nil
}

func (r *Repository) GetSingleMassCalculation(id int) (ds.MassCalculation, error) {
	if id < 0 {
		return ds.MassCalculation{}, errors.New("invalid id, it must be >= 0")
	}
	var calculation ds.MassCalculation
	err := r.db.Where("id = ?", id).First(&calculation).Error
	if err != nil {
		return ds.MassCalculation{}, ErrorNotFound
	}

	return calculation, nil
}

func (r *Repository) FormMassCalculation(id int, status string) (ds.MassCalculation, error) {
	calculation, err := r.GetSingleMassCalculation(id)
	if err != nil {
		return ds.MassCalculation{}, err
	}

	if calculation.Status != "draft" {
		return ds.MassCalculation{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&calculation).Updates(ds.MassCalculation{
		Status: status,
		DateForm: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.MassCalculation{}, err
	}

	return calculation, nil
}

func (r *Repository) ModerateMassCalculation(id int, status string) (ds.MassCalculation, error) {
	if status != "completed" && status != "rejected" {
		return ds.MassCalculation{}, errors.New("wrong status")
	}

	userId := r.GetUserID()

	calculation, err := r.GetSingleMassCalculation(id)
	if err != nil {
		return ds.MassCalculation{}, err
	} else if calculation.Status != "formed" {
		return ds.MassCalculation{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&calculation).Updates(ds.MassCalculation{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: sql.NullInt64{
			Int64: int64(userId),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.MassCalculation{}, err
	}

	if status == "completed" {
		reactionCalculations, err := r.GetReactionCalculations(calculation.ID)
		if err != nil {
			return ds.MassCalculation{}, err
		}
		for _, reactionCalculation := range reactionCalculations {
			reaction, err := r.GetReaction(reactionCalculation.ReactionID)
			if err != nil {
				return ds.MassCalculation{}, err
			}
			mass, err := CalculateMass(reactionCalculation.OutputMass, reaction.ConversationFactor, calculation.OutputKoef)
			if err != nil {
				return ds.MassCalculation{}, err
			}
			err = r.db.Model(&reactionCalculation).Updates(ds.ReactionCalculation{
				InputMass: mass,
			}).Error
			if err != nil {
				return ds.MassCalculation{}, err
			}
		}
	}

	return calculation, nil
}
