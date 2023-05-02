package service

import (
	"context"
	"time"

	"github.com/SultanKs4/wassistant/internal/entity"
)

type messageService struct {
	messageRepo entity.MessageRepository
	contactRepo entity.ContactRepository
}

func NewMessageService(msgRepo entity.MessageRepository, ctcRepo entity.ContactRepository) entity.MessageService {
	return &messageService{messageRepo: msgRepo, contactRepo: ctcRepo}
}

func (s messageService) Store(ctx context.Context, m *entity.Message, ctc *entity.Contact) (err error) {
	ctxTo, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	dbContact, err := s.contactRepo.FirstOrCreate(ctxTo, ctc)
	if err != nil {
		return
	}
	m.ContactID = dbContact.ID
	err = s.messageRepo.Store(ctxTo, m)
	return
}
