package repository

import (
	minioInclude "Backend/internal/app/minio"

	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	mc     *minio.Client
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	mc, err := minioInclude.InitMinio()
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
		mc: mc,
	}, nil
}

func (r *Repository) GetUserID() int {
	return 1
}
