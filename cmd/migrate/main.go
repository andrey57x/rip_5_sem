package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"Backend/internal/app/ds"
	"Backend/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Reaction{},
		&ds.Calculation{},
		&ds.ReactionCalculation{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}