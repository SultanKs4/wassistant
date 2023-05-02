package cli

import (
	"context"

	"go.mau.fi/whatsmeow/types"

	"github.com/SultanKs4/wassistant/internal/entity"
)

type MessageHandler struct {
	msgService entity.MessageService
}

func NewMessageHandler(msgService entity.MessageService) *MessageHandler {
	return &MessageHandler{msgService: msgService}
}

func (h MessageHandler) StoreMessageWhatsapp(jidClient *types.JID, info types.MessageInfo, contactName types.ContactInfo, message string) (err error) {
	jidReceiver := info.Chat.User
	if jidReceiver == info.Sender.User && jidReceiver != jidClient.User {
		jidReceiver = jidClient.User
	}
	c := &entity.Contact{
		Jid:      info.Sender.User,
		PushName: info.PushName,
		FullName: contactName.FullName,
	}
	m := &entity.Message{
		JidReceiver: jidReceiver,
		Message:     message,
		Timestamp:   info.Timestamp,
		IsGroup:     info.IsGroup,
	}
	err = h.msgService.Store(context.Background(), m, c)
	return err
}
