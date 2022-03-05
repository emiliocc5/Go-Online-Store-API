package repository

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	IDBClient interface {
		GetClient() (*gorm.DB, error)
	}
	DBClient struct{}
)

func (dbc *DBClient) GetClient() (*gorm.DB, error) {
	return connectDatabase()
}

func connectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=goApiUser password=1234 dbname=OnlineStore port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		//TODO LOG ERROR
		return nil, err
	}

	err1 := database.AutoMigrate(&models.Cart{})
	if err1 != nil {
		//TODO LOG ERROR
		return nil, err1
	}

	return database, nil
}
