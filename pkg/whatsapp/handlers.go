package whatsapp

import (
	"fmt"
	"strings"

	"github.com/SultanKs4/wassistant/message/delivery/cli"
	"github.com/SultanKs4/wassistant/message/repository"
	"github.com/SultanKs4/wassistant/message/service"
	"github.com/SultanKs4/wassistant/pkg/db"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func logMessageToDb(info types.MessageInfo, message string) error {
	gdb, err := db.CreateDb()
	if err != nil {
		return fmt.Errorf("create db error: %v", err.Error())
	}
	pg := db.NewPg(gdb)
	defer pg.Disconnect()
	mRepo := repository.NewMessageRepository(gdb)
	mServ := service.NewMessageService(mRepo)
	mHand := cli.NewMessageHandler(mServ)

	contactName, err := client.Store.Contacts.GetContact(info.Chat)
	if err != nil {
		return fmt.Errorf("get name contact: %v", err.Error())
	}

	err = mHand.StoreMessageWhatsapp(client.Store.ID, info, contactName, message)
	if err != nil {
		return fmt.Errorf("store message to db error: %v", err.Error())
	}
	return nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		message := strings.ToLower(v.Message.GetConversation())

		if v.Info.Chat.Server != "broadcast" && message != "" {
			err := clientIsOk()
			if err != nil {
				fmt.Println("error client: ", err.Error())
			}

			err = logMessageToDb(v.Info, message)
			if err != nil {
				fmt.Println("error log message to db: ", err.Error())
			}

			private := true
			if v.Info.IsFromMe || v.Info.IsGroup {
				private = false
			}

			if private {
				switch message {
				case "cat":
					time, err := cat(v.Info.Chat.User, v.Info.PushName)
					if err != nil {
						fmt.Println("error cat actions: ", err.Error())
					}
					fmt.Println("message sended at: ", time)
				}
			}

		}

	}
}
