package whatsapp

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func sendText(rjid types.JID, message string) (string, error) {
	// Set Chat Presence
	composeStatus(rjid, true, false)
	defer composeStatus(rjid, false, false)

	msgID := whatsmeow.GenerateMessageID()
	msgContent := &waProto.Message{
		Conversation: proto.String(message),
	}

	sendTime, err := client.SendMessage(rjid, msgID, msgContent)
	if err != nil {
		return "", err
	}
	return sendTime.String(), nil
}
