package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/SultanKs4/wassistant/internal/entity"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) entity.MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Store(c context.Context, m *entity.Message) error {
	res := r.db.WithContext(c).Create(&m)
	return res.Error
}
