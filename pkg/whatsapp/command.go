package whatsapp

import (
	"fmt"
)

func cat(phone string, pushname string) (string, error) {
	rjid := getJid(phone)
	if rjid.IsEmpty() {
		return "", fmt.Errorf("phone not registered")
	}

	msgSend := fmt.Sprintf("catto lul %v 😺", pushname)

	time, err := sendText(rjid, msgSend)
	if err != nil {
		return "", fmt.Errorf("error send message: %v", err.Error())
	}
	return time, nil
}
