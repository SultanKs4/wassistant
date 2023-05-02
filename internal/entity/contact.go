package entity

import (
	"context"
)

type Contact struct {
	ID       uint   `gorm:"primarykey"`
	Jid      string `gorm:"uniques"`
	PushName string
	FullName string
	Messages []Message
}

type contact interface {
	FirstOrCreate(ctx context.Context, contact *Contact) (Contact, error)
}

type ContactRepository interface {
	contact
}

type contactService interface {
	contact
}
