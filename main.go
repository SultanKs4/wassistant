package main

import (
	"os"
	"os/signal"
	"syscall"

	database "github.com/SultanKs4/wassistant/pkg/db"
	"github.com/SultanKs4/wassistant/pkg/whatsapp"
)

func main() {
	db, err := database.MigrateDbPg()
	if err != nil {
		panic(err)
	}
	pg := database.NewPg(db)
	defer pg.Disconnect()

	err = whatsapp.Connect()
	if err != nil {
		panic(err)
	}
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	defer whatsapp.Disconnect()
}
