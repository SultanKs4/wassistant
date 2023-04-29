package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type waCli struct {
	client         *whatsmeow.Client
	eventHandlerID uint32
}

func Connect(db *sql.DB) (*waCli, error) {
	container, err := newContainer(db)
	if err != nil {
		return nil, err
	}

	deviceStore, err := newDeviceStore(container)
	if err != nil {
		return nil, err
	}

	wa := newClient(deviceStore)

	if err = wa.login(); err != nil {
		return nil, err
	}

	wa.register()

	resp, err := wa.client.CheckUpdate()
	if err != nil {
		return nil, err
	}
	fmt.Println("connected to whatsapp version: ", resp.CurrentVersion)

	return wa, nil
}

func newContainer(db *sql.DB) (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(db, "pgx", dbLog)
	if err := container.Upgrade(); err != nil {
		return nil, err
	}

	return container, nil
}

func newDeviceStore(container *sqlstore.Container) (*store.Device, error) {
	// device default value
	store.DeviceProps.Os = proto.String("wassistant")
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	return deviceStore, nil
}

func newClient(deviceStore *store.Device) *waCli {
	clientLog := waLog.Stdout("Client", "INFO", true)
	return &waCli{client: whatsmeow.NewClient(deviceStore, clientLog)}
}

func (wa *waCli) login() error {
	if wa.client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := wa.client.GetQRChannel(context.Background())
		err := wa.client.Connect()
		if err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				return fmt.Errorf("Login event: %v, error: %v", evt.Event, evt.Error)
			}
		}
	} else {
		// Already logged in, just connect
		err := wa.client.Connect()
		if err != nil {
			return err
		}
	}

	wa.client.SendPresence(types.PresenceAvailable)
	wa.client.SendPresence(types.PresenceUnavailable)

	return nil
}

func (wa *waCli) register() {
	wa.eventHandlerID = wa.client.AddEventHandler(wa.MessageHandlerWhatsapp)
}

func (wa *waCli) Disconnect() {
	wa.client.Disconnect()
}
