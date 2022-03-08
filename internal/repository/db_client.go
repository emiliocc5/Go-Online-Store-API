package repository

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sync"
)

type (
	IDBClient interface {
		GetClient() (*gorm.DB, error)
	}
	DBClientImpl struct{}
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

func GetDbClient() *DBClientImpl {
	return &DBClientImpl{}
}

func (dbc *DBClientImpl) GetClient() (*gorm.DB, error) {
	return getClientInstance()
}

func getClientInstance() (*gorm.DB, error) {
	once.Do(func() {
		dbInstance, errorInstance = connectDatabase()
	})
	return dbInstance, errorInstance
}

func connectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost port=5432 user=goApiUser password=1234 dbname=OnlineStore sslmode=disable")
	//dsn := fmt.Sprintf("host=%+v user=%+v password=%+v dbname=%+v port=%+v sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)

	fmt.Println("Openning connection")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		//TODO LOG ERROR
		fmt.Println(fmt.Sprintf("error message: %s", err))
		return nil, err
	}
	err1 := migrateTables(db)
	if err1 != nil {
		return nil, err1
	}

	return db, nil
}

func migrateTables(db *gorm.DB) error {
	err1 := db.AutoMigrate(&models.Cart{})
	if err1 != nil {
		//TODO LOG ERROR
		return err1
	}
	err2 := db.AutoMigrate(&models.Product{})
	if err2 != nil {
		//TODO LOG ERROR
		return err2
	}
	err3 := db.AutoMigrate(&models.ProductCart{})
	if err3 != nil {
		return err3
	}
	err4 := db.AutoMigrate(&models.Client{})
	if err4 != nil {
		return err4
	}

	return nil
}
