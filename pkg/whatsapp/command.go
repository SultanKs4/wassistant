package whatsapp

import (
	"fmt"

	database "github.com/SultanKs4/wassistant/pkg/db"
	"github.com/SultanKs4/wassistant/repository"
)

func cat(phone string) (string, error) {
	rjid := getJid(phone)
	if rjid.IsEmpty() {
		return "", fmt.Errorf("phone not registered")
	}

	db, err := database.CreateDb()
	if err != nil {
		panic(err)
	}
	pg := &database.Pg{Db: db}
	defer pg.Disconnect()
	mMessage := repository.NewMessageRepository(db)
	defer mMessage.GetDataMessage()

	msg := mMessage.GetDataMessage()
	if msg == nil {
		return "", fmt.Errorf("data from db message empty")
	}

	msgSend := fmt.Sprintf("catto lul 😺: %v", msg.Message)

	time, err := sendText(rjid, msgSend)
	if err != nil {
		return "", fmt.Errorf("error send message: %v", err.Error())
	}
	return time, nil
}
