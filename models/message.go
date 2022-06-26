package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Jid     string
	Name    string
	Message string
}
