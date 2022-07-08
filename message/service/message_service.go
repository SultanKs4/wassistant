package service

import (
	"context"
	"time"

	"github.com/SultanKs4/wassistant/message/entity"
)

type messageService struct {
	messageRepo entity.MessageRepository
}

func NewMessageService(messageRepo entity.MessageRepository) entity.MessageService {
	return &messageService{messageRepo}
}

func (s messageService) Store(c context.Context, m *entity.Message) (err error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	err = s.messageRepo.Store(ctx, m)
	return
}
