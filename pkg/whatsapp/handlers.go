package whatsapp

import (
	"fmt"
	"strings"

	_cRepo "github.com/SultanKs4/wassistant/contact/repository"
	"github.com/SultanKs4/wassistant/message/delivery/cli"
	_mRepo "github.com/SultanKs4/wassistant/message/repository"
	"github.com/SultanKs4/wassistant/message/service"
	"github.com/SultanKs4/wassistant/pkg/db"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func logMessageToDb(info types.MessageInfo, contactName types.ContactInfo, message string, ch chan error) {
	gdb, err := db.CreateDbPg()
	if err != nil {
		ch <- fmt.Errorf("create db error: %v", err.Error())
	}
	pg := db.NewPg(gdb)
	defer pg.Disconnect()
	mRepo := _mRepo.NewMessageRepository(gdb)
	cRepo := _cRepo.NewContactRepository(gdb)
	mServ := service.NewMessageService(mRepo, cRepo)
	mHand := cli.NewMessageHandler(mServ)

	err = mHand.StoreMessageWhatsapp(client.Store.ID, info, contactName, message)
	if err != nil {
		ch <- fmt.Errorf("store message to db error: %v", err.Error())
	}
	ch <- nil
}

func checkMessage(v *events.Message) (message string) {
	message = v.Message.GetConversation()
	if message == "" {
		switch v.Info.MediaType {
		case "vcard":
			message = fmt.Sprintf("contact: %v", v.Message.GetContactMessage().GetDisplayName())
		case "location":
			lat := v.Message.GetLocationMessage().GetDegreesLatitude()
			long := v.Message.GetLocationMessage().GetDegreesLongitude()
			message = fmt.Sprintf("location: (%f, %f)", lat, long)
		default:
			message = v.Message.GetExtendedTextMessage().GetText()
		}
	}
	return
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.Chat.Server == "broadcast" {
			break
		}

		err := clientIsOk()
		if err != nil {
			fmt.Println("error client: ", err.Error())
		}

		message := checkMessage(v)

		if message != "" {
			contactName, err := client.Store.Contacts.GetContact(v.Info.Sender)
			if err != nil {
				fmt.Println("get name contact: ", err.Error())
			}

			var errCh = make(chan error, 0)
			go logMessageToDb(v.Info, contactName, message, errCh)
			err = <-errCh
			if err != nil {
				fmt.Println("error log message to db: ", err.Error())
			}
			close(errCh)

			private := true
			if v.Info.IsFromMe || v.Info.IsGroup {
				private = false
			}

			if private {
				switch strings.ToLower(message) {
				case "cat":
					name := v.Info.PushName
					if contactName.FullName != "" {
						name = contactName.FullName
					}
					err := cat(v.Info.Chat.User, name)
					if err != nil {
						fmt.Println("error cat actions: ", err.Error())
					}
				}
			}
		}
	}
}
