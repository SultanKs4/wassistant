package entity

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	JidSender   string
	PushName    string
	FullName    string
	JidReceiver string
	Message     string
	Timestamp   time.Time
	IsGroup     bool
}

type message interface {
	Store(c context.Context, m *Message) error
}

type MessageRepository interface {
	message
}

type MessageService interface {
	message
}
