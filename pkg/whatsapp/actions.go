package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func sendText(rjid types.JID, message string) error {
	// Set Chat Presence
	composeStatus(rjid, true, false)
	defer composeStatus(rjid, false, false)

	msgID := whatsmeow.GenerateMessageID()
	msgContent := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err := client.SendMessage(context.Background(), rjid, msgID, msgContent)
	if err != nil {
		return err
	}
	return nil
}
