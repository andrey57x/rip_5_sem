package apitypes

import (
	"Backend/internal/app/ds"
	"database/sql"
	"time"
)

type CalculationJSON struct {
	ID             int          `json:"id"`
	OutputKoef     float32      `json:"output_koef"`
	Status         string       `json:"status"`
	DateCreate     time.Time    `json:"date_create"`
	DateForm       sql.NullTime `json:"date_form"`
	DateFinish     sql.NullTime `json:"date_finish"`
	CreatorLogin   string       `json:"creator_login"`
	ModeratorLogin string       `json:"moderator_login"`
}

func CalculationToJSON(c ds.Calculation, creatorLogin, moderatorLogin string) CalculationJSON {
	return CalculationJSON{
		ID:             c.ID,
		OutputKoef:     c.OutputKoef,
		Status:         c.Status,
		DateCreate:     c.DateCreate,
		DateForm:       c.DateForm,
		DateFinish:     c.DateFinish,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: moderatorLogin,
	}
}

func CalculationFromJSON(c CalculationJSON) ds.Calculation {
	return ds.Calculation{
		OutputKoef: c.OutputKoef,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
