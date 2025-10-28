package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrorNotAllowed = errors.New("you are not allowed")
	ErrorNoDraft    = errors.New("no draft for this user")
	ErrorNotFound   = errors.New("not found")
	ErrorDeleted    = errors.New("calculations is deleted")
)

func (r *Repository) GetMassCalculationReactions(id int) ([]ds.Reaction, ds.MassCalculation, error) {
	calculation, err := r.GetSingleMassCalculation(id)
	if err != nil {
		return []ds.Reaction{}, ds.MassCalculation{}, err
	}
	if calculation.Status == "deleted" {
		return []ds.Reaction{}, ds.MassCalculation{}, ErrorDeleted
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
func (r *Repository) GetCartCount(creatorID uuid.UUID) int {
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

func (r *Repository) CheckCurrentMassCalculationDraft(creatorID uuid.UUID) (ds.MassCalculation, error) {
	var calculation ds.MassCalculation

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&calculation)
	if res.Error != nil {
		return ds.MassCalculation{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.MassCalculation{}, ErrorNoDraft
	}
	return calculation, nil
}

func (r *Repository) GetMassCalculationDraft(creatorID uuid.UUID) (ds.MassCalculation, bool, error) {
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
	sub := r.db.Where("status != 'deleted'")
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

func (r *Repository) ModerateMassCalculation(id int, status string, currUserId uuid.UUID) (ds.MassCalculation, error) {
	if status != "completed" && status != "rejected" {
		return ds.MassCalculation{}, errors.New("wrong status")
	}

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
		ModeratorID: uuid.NullUUID{
			UUID:  currUserId,
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.MassCalculation{}, err
	}

	return calculation, nil
}

func (r *Repository) UpdateReactionCalculationResult(calculationID, reactionID int, inputMass float32) error {
	result := r.db.Model(&ds.ReactionCalculation{}).
		Where("calculation_id = ? AND reaction_id = ?", calculationID, reactionID).
		Update("input_mass", inputMass)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrorNotFound
	}
	return nil
}

func (r *Repository) GetCompletedReactionsCount(calculationID int) (int, error) {
	var count int64
	err := r.db.Model(&ds.ReactionCalculation{}).
		Where("calculation_id = ? AND input_mass IS NOT NULL AND input_mass != 0", calculationID).
		Count(&count).Error
	
	return int(count), err
}
