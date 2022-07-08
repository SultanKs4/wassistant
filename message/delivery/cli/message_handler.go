package cli

import (
	"context"

	"go.mau.fi/whatsmeow/types"

	"github.com/SultanKs4/wassistant/message/entity"
)

type messageHandler struct {
	msgService entity.MessageService
}

func NewMessageHandler(msgService entity.MessageService) *messageHandler {
	return &messageHandler{msgService: msgService}
}

func (h messageHandler) StoreMessageWhatsapp(jidClient *types.JID, info types.MessageInfo, contactName types.ContactInfo, message string) (err error) {
	jidReceiver := info.Chat.User
	if jidReceiver == info.Sender.User && jidReceiver != jidClient.User {
		jidReceiver = jidClient.User
	}
	m := &entity.Message{
		JidSender:   info.Sender.User,
		PushName:    info.PushName,
		FullName:    contactName.FullName,
		JidReceiver: jidReceiver,
		Message:     message,
		Timestamp:   info.Timestamp,
		IsGroup:     info.IsGroup,
	}
	err = h.msgService.Store(context.Background(), m)
	return err
}
