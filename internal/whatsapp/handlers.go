package whatsapp

import (
	"context"
	"fmt"

	"github.com/SultanKs4/wassistant/pkg/utils"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func checkMessage(v *events.Message) (message string) {
	message = v.Message.GetConversation()
	if message != "" {
		return
	}

	switch v.Info.MediaType {
	case "document":
	case "ptt":
	case "video":
	case "gif":
	case "image":
		// img := v.Message.GetImageMessage()
	case "vcard":
		message = fmt.Sprintf("contact: %v", v.Message.GetContactMessage().GetDisplayName())
	case "contact_array":

	case "location":
		loc := v.Message.GetLocationMessage()
		lat := loc.GetDegreesLatitude()
		long := loc.GetDegreesLongitude()
		message = fmt.Sprintf("location: (%f, %f)", lat, long)
	case "product":
		// prod := v.Message.GetProductMessage()
	default:
		message = v.Message.GetExtendedTextMessage().GetText()
	}
	return
}

func (wa *waCli) handler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected, *events.PushNameSetting:
		if len(wa.client.Store.PushName) == 0 {
			return
		}
		if err := wa.client.SendPresence(types.PresenceAvailable); err != nil {
			wa.client.Log.Warnf("Failed to send available presence: %v", err)
		}
	// case *events.Picture:

	case *events.Message:
		if v.Info.Chat.Server == "broadcast" {
			break
		}

		if err := wa.checkClient(); err != nil {
			wa.client.Log.Errorf("error client: %w", err)
		}

		contactName, err := wa.getContact(v.Info.Sender)
		if err != nil {
			fmt.Println("get name contact: ", err.Error())
		}

		message := checkMessage(v)
		if message == "" {
			break
		}

		if err := wa.msgHandler.StoreMessageWhatsapp(wa.client.Store.ID, v.Info, contactName, message); err != nil {
			fmt.Println("store message to db error: ", err.Error())
		}

		private := true
		if v.Info.IsFromMe || v.Info.IsGroup {
			private = false
		}

		if !private {
			break
		}

		switch utils.Sanitize(message) {
		case "cat":
			name := v.Info.PushName
			if contactName.FullName != "" {
				name = contactName.FullName
			}
			err := wa.cmdCat(context.Background(), v.Info.Chat, name)
			if err != nil {
				fmt.Println("error cat actions: ", err.Error())
			}
		}
	}
}
