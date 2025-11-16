package migrate

import (
	"Backend/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Reaction{},
		&ds.MassCalculation{},
		&ds.ReactionCalculation{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
