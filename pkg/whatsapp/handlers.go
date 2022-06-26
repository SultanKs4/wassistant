package whatsapp

import (
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types/events"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		message := strings.ToLower(v.Message.GetConversation())
		server := v.Info.Chat.Server
		private := true
		if v.Info.IsFromMe {
			private = false
		}
		if v.Info.IsGroup {
			private = false
		}
		if private && server != "broadcast" && message != "" {
			err := clientIsOk()
			if err != nil {
				fmt.Println("error client: ", err)
			}

			err = client.MarkRead([]string{v.Info.ID}, time.Now(), v.Info.Chat, *client.Store.ID)
			if err != nil {
				fmt.Println("error mark read: ", err)
			}

			fmt.Println("JID Chat: ", v.Info.Chat)
			fmt.Println("timestamps: ", v.Info.Timestamp)
			fmt.Println("sender: ", v.Info.Sender.User)
			fmt.Println("info pushname: ", v.Info.PushName)
			fmt.Println("Received a message: ", message)
			switch message {
			case "cat":
				time, err := cat(v.Info.Chat.User)
				if err != nil {
					fmt.Println("error cat actions: ", err)
				}
				fmt.Println("message sended at: ", time)
			}
		}

	}
}
