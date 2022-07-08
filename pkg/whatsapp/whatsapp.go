package whatsapp

import (
	"context"
	"errors"
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

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var client *whatsmeow.Client

func initContainer(databaseUrl string) (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("pgx", databaseUrl, dbLog)

	if err != nil {
		return nil, err
	}

	return container, nil
}

func initDeviceStore(container sqlstore.Container) (*store.Device, error) {
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	store.DeviceProps.Os = proto.String("wassistant")
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	return deviceStore, nil
}

func getJid(phone string) types.JID {
	phones := []string{}
	phones = append(phones, "+"+phone)
	jidInfos, err := client.IsOnWhatsApp(phones)
	// If Phone Number is Registered as JID
	if err == nil && jidInfos[0].IsIn {
		// Return JID Information
		return jidInfos[0].JID
	}
	return types.EmptyJID
}

func clientIsOk() error {
	// Make Sure WhatsApp Client is Connected
	if !client.IsConnected() {
		return errors.New("WhatsApp Client is not Connected")
	}

	// Make Sure WhatsApp Client is Logged In
	if !client.IsLoggedIn() {
		return errors.New("WhatsApp Client is not Logged In")
	}

	return nil
}

func composeStatus(rjid types.JID, isComposing bool, isAudio bool) {
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
	_ = client.SendChatPresence(rjid, typeCompose, typeComposeMedia)
}

func initClient(deviceStore *store.Device) {
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)
}

func login() error {
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := client.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

func Connect() error {
	container, err := initContainer(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	deviceStore, err := initDeviceStore(*container)
	if err != nil {
		return err
	}

	initClient(deviceStore)

	err = login()
	if err != nil {
		return err
	}

	_ = client.SendPresence(types.PresenceAvailable)
	_ = client.SendPresence(types.PresenceUnavailable)

	resp, err := client.CheckUpdate()
	if err != nil {
		return err
	}
	fmt.Println("connected to whatsapp version: ", resp.CurrentVersion)

	return nil
}

func Disconnect() {
	client.Disconnect()
}
