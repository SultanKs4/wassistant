package whatsapp

import (
	"context"
	"errors"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (wa *waCli) composeStatus(rjid types.JID, isComposing bool, isAudio bool) {
	// Set Compose Status
	var typeCompose types.ChatPresence
	if isComposing {
		typeCompose = types.ChatPresenceComposing
	} else {
		typeCompose = types.ChatPresencePaused
	}

	// Set Compose Media Audio (Recording) or Text (Typing)
	var typeComposeMedia types.ChatPresenceMedia
	if isAudio {
		typeComposeMedia = types.ChatPresenceMediaAudio
	} else {
		typeComposeMedia = types.ChatPresenceMediaText
	}

	// Send Chat Compose Status
	_ = wa.client.SendChatPresence(rjid, typeCompose, typeComposeMedia)
}

func (wa *waCli) SendText(rjid types.JID, message string) error {
	// Set Chat Presence
	wa.composeStatus(rjid, true, false)
	defer wa.composeStatus(rjid, false, false)

	_, err := wa.client.SendMessage(context.Background(), rjid, &waProto.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		return err
	}
	return nil
}

func (wa *waCli) getJid(phone string) types.JID {
	phones := []string{}
	phones = append(phones, "+"+phone)
	jidInfos, err := wa.client.IsOnWhatsApp(phones)
	// If Phone Number is Registered as JID
	if err == nil && jidInfos[0].IsIn {
		// Return JID Information
		return jidInfos[0].JID
	}
	return types.EmptyJID
}

func (wa *waCli) checkClient() error {
	// Make Sure WhatsApp Client is Connected
	if !wa.client.IsConnected() {
		return errors.New("WhatsApp Client is not Connected")
	}

	// Make Sure WhatsApp Client is Logged In
	if !wa.client.IsLoggedIn() {
		return errors.New("WhatsApp Client is not Logged In")
	}

	return nil
}

func (wa *waCli) getContact(jid types.JID) (contactName types.ContactInfo, err error) {
	contactName, err = wa.client.Store.Contacts.GetContact(jid)
	return
}
