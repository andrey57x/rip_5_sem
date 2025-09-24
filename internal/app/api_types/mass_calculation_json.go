package apitypes

import (
	"Backend/internal/app/ds"
	"time"
)

// MassCalculationJSON model
// @Description Model for calculation
// @Tags calculations
type MassCalculationJSON struct {
	ID             int        `json:"id"`
	OutputKoef     *float32   `json:"output_koef"`
	Status         string     `json:"status"`
	DateCreate     time.Time  `json:"date_create"`
	DateForm       *time.Time `json:"date_form"`
	DateFinish     *time.Time `json:"date_finish"`
	CreatorLogin   string     `json:"creator_login"`
	ModeratorLogin *string    `json:"moderator_login"`
}

func MassCalculationToJSON(c ds.MassCalculation, creatorLogin, moderatorLogin string) MassCalculationJSON {
	var dateForm, dateFinish *time.Time
	if c.DateForm.Valid {
		dateForm = &c.DateForm.Time
	}

	if c.DateFinish.Valid {
		dateFinish = &c.DateFinish.Time
	}

	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}

	var outputKoef *float32
	if c.OutputKoef != 0 {
		outputKoef = &c.OutputKoef
	}

	return MassCalculationJSON{
		ID:             c.ID,
		OutputKoef:     outputKoef,
		Status:         c.Status,
		DateCreate:     c.DateCreate,
		DateForm:       dateForm,
		DateFinish:     dateFinish,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: mLogin,
	}
}

func MassCalculationFromJSON(c MassCalculationJSON) ds.MassCalculation {
	if c.OutputKoef == nil {
		return ds.MassCalculation{}
	}
	return ds.MassCalculation{
		OutputKoef: *c.OutputKoef,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
