package whatsapp

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/types"
)

func (wa *waCli) cmdCat(ctx context.Context, rjid types.JID, name string) error {
	msgSend := fmt.Sprintf("catto lul %v 😺", name)
	ctxTo, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err := wa.SendText(ctxTo, rjid, msgSend)
	if err != nil {
		return fmt.Errorf("error send message: %w", err)
	}
	return nil
}
