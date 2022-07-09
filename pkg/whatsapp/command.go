package whatsapp

import (
	"fmt"
)

func cat(phone string, name string) error {
	rjid := getJid(phone)
	if rjid.IsEmpty() {
		return fmt.Errorf("phone not registered")
	}

	msgSend := fmt.Sprintf("catto lul %v 😺", name)

	err := sendText(rjid, msgSend)
	if err != nil {
		return fmt.Errorf("error send message: %v", err.Error())
	}
	return nil
}
