package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ctcRepo "github.com/SultanKs4/wassistant/internal/contact/repository"
	_msgHand "github.com/SultanKs4/wassistant/internal/message/delivery/cli"
	_msgRepo "github.com/SultanKs4/wassistant/internal/message/repository"
	_msgServ "github.com/SultanKs4/wassistant/internal/message/service"
	"github.com/SultanKs4/wassistant/pkg/db/postgresql"
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
	msgHandler     *_msgHand.MessageHandler
}

func newContainer(db *sql.DB) (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(db, "pgx", dbLog)
	if err := container.Upgrade(); err != nil {
		return nil, err
	}

	return container, nil
}

func newDeviceStore(container *sqlstore.Container, jid types.JID) (*store.Device, error) {
	var (
		deviceStore *store.Device
		err         error
	)
	// device default value
	store.DeviceProps.Os = proto.String("wassistant")
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	if jid.IsEmpty() {
		deviceStore, err = container.GetFirstDevice()
		if err != nil {
			return nil, err
		}
	} else {
		deviceStore, err = container.GetDevice(jid)
		if err != nil {
			return nil, err
		}
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
				if evt.Error != nil {
					return fmt.Errorf("Login event: %v, error: %v", evt.Event, evt.Error)
				} else {
					return nil
				}
			}
		}
	} else {
		// Already logged in, just connect
		err := wa.client.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (wa *waCli) register() {
	wa.eventHandlerID = wa.client.AddEventHandler(wa.handler)
}

func (wa *waCli) Disconnect() {
	wa.client.Disconnect()
}

func Connect(db *postgresql.Postgre) (*waCli, error) {
	container, err := newContainer(db.GetSqlDb())
	if err != nil {
		return nil, err
	}

	deviceStore, err := newDeviceStore(container, types.EmptyJID)
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

	mRepo := _msgRepo.NewMessageRepository(db.Db)
	cRepo := _ctcRepo.NewContactRepository(db.Db)
	mServ := _msgServ.NewMessageService(mRepo, cRepo)
	wa.msgHandler = _msgHand.NewMessageHandler(mServ)

	return wa, nil
}
