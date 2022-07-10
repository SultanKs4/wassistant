package entity

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	JidReceiver string
	Message     string
	Timestamp   time.Time
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
