package whatsapp

import (
	"fmt"

	"go.mau.fi/whatsmeow/types"
)

func cat(waCli *waCli, rjid types.JID, name string) error {
	msgSend := fmt.Sprintf("catto lul %v 😺", name)
	err := waCli.SendText(rjid, msgSend)
	if err != nil {
		return fmt.Errorf("error send message: %v", err.Error())
	}
	return nil
}
