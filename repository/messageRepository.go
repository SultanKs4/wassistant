package repository

import (
	"errors"

	"github.com/SultanKs4/wassistant/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	Db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (m *MessageRepository) GetDataMessage() *models.Message {
	var message models.Message
	err := m.Db.First(&message).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &message
}
