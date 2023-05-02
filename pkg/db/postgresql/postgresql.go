package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/SultanKs4/wassistant/config"
	"github.com/SultanKs4/wassistant/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgre struct {
	Db *gorm.DB
}

func NewPg(cfg config.PostgresConfig) (*Postgre, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed connect database gorm: %v", err.Error())
	}

	return &Postgre{db}, nil
}

func (pg *Postgre) Migrate() error {
	return pg.Db.AutoMigrate(&entity.Contact{}, &entity.Message{})
}

func (pg *Postgre) GetSqlDb() *sql.DB {
	db, _ := pg.Db.DB()
	return db
}

func (pg *Postgre) Disconnect() error {
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
