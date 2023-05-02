package entity

import (
	"context"
	"time"
)

type Message struct {
	ID          uint `gorm:"primarykey"`
	JidReceiver string
	Message     string
	Timestamp   time.Time `gorm:"index"`
	IsGroup     bool
	ContactID   uint
}

type message interface {
	Store(c context.Context, m *Message) error
}

type MessageRepository interface {
	message
}

type MessageService interface {
	Store(ctx context.Context, message *Message, contact *Contact) error
}
