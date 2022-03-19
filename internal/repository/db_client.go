package repository

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	once          sync.Once
	dbInstance    *gorm.DB
	errorInstance error

	DbHost     = os.Getenv("DB_HOST")
	DbUser     = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName     = os.Getenv("DB_NAME")
	DbPort     = os.Getenv("DB_PORT")
)

func init() {
	logger = utils.GetLogger()
}

func GetClient() (*gorm.DB, error) {
	return getClientInstance()
}

func getClientInstance() (*gorm.DB, error) {
	once.Do(func() {
		dbInstance, errorInstance = connectDatabase()
	})
	return dbInstance, errorInstance
}

//TODO get variables from env
func connectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost port=5432 user=goApiUser password=1234 dbname=OnlineStore sslmode=disable")
	//dsn := fmt.Sprintf("host=%+v user=%+v password=%+v dbname=%+v port=%+v sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)

	logger.Info("Opening connections")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("error message: %s", err)
		return nil, err
	}
	err1 := migrateTables(db)
	if err1 != nil {
		return nil, err1
	}

	return db, nil
}

func migrateTables(db *gorm.DB) error {
	err1 := db.AutoMigrate(&models.Cart{},
		&models.Product{},
		&models.ProductCart{},
		&models.Client{})
	if err1 != nil {
		return err1
	}

	return nil
}
