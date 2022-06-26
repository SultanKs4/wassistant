package db

import (
	"fmt"
	"os"

	"github.com/SultanKs4/wassistant/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Pg struct {
	Db *gorm.DB
}

func NewPg(db *gorm.DB) *Pg {
	return &Pg{Db: db}
}

func CreateDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
	if err != nil {
		return db, fmt.Errorf("failed connect database gorm: %v", err.Error())
	}

	return db, nil
}

func MigrateDbPg() (*gorm.DB, error) {
	db, err := CreateDb()
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Message{})

	return db, nil
}

func (pg *Pg) Disconnect() error {
	sqlDb, err := pg.Db.DB()
	if err != nil {
		return err
	}
	err = sqlDb.Close()
	if err != nil {
		return err
	}
	return nil
}
