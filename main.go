package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/SultanKs4/wassistant/pkg/db"
	"github.com/SultanKs4/wassistant/pkg/whatsapp"
)

func main() {
	pgdb, err := db.MigrateDbPg()
	if err != nil {
		panic(err)
	}
	pg := db.NewPg(pgdb)
	pg.Disconnect()

	rdb, err := db.CreateDbRedis()
	if err != nil {
		panic(err)
	}
	redb := db.NewRedis(rdb)
	redb.Rdb.Close()

	err = whatsapp.Connect()
	if err != nil {
		panic(err)
	}
	defer whatsapp.Disconnect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
