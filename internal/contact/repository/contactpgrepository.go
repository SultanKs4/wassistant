package repository

import (
	"context"

	"github.com/SultanKs4/wassistant/internal/entity"
	"gorm.io/gorm"
)

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) entity.ContactRepository {
	return &contactRepository{db}
}

func (r *contactRepository) FirstOrCreate(ctx context.Context, contact *entity.Contact) (entity.Contact, error) {
	var resCon *entity.Contact
	res := r.db.WithContext(ctx).FirstOrCreate(&resCon, &contact)
	return *resCon, res.Error
}
