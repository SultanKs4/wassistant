package service

import (
	"context"
	"time"

	"github.com/SultanKs4/wassistant/entity"
)

type messageService struct {
	messageRepo entity.MessageRepository
	contactRepo entity.ContactRepository
}

func NewMessageService(msgRepo entity.MessageRepository, ctcRepo entity.ContactRepository) entity.MessageService {
	return &messageService{messageRepo: msgRepo, contactRepo: ctcRepo}
}

func (s messageService) Store(c context.Context, m *entity.Message, contact *entity.Contact) (err error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	dbContact, err := s.contactRepo.FirstOrCreate(ctx, contact)
	if err != nil {
		return
	}
	m.ContactID = dbContact.ID
	err = s.messageRepo.Store(ctx, m)
	return
}
