package db

import (
	"fmt"

	"github.com/SultanKs4/wassistant/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgre struct {
	*gorm.DB
}

func NewPg(dsn string) (*Postgre, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed connect database gorm: %v", err.Error())
	}

	return &Postgre{db}, nil
}

func (pg *Postgre) Migrate() {
	pg.AutoMigrate(&entity.Contact{}, &entity.Message{})
}

func (pg *Postgre) Disconnect() error {
	sqlDb, err := pg.DB.DB()
	if err != nil {
		return err
	}
	err = sqlDb.Close()
	if err != nil {
		return err
	}
	return nil
}
