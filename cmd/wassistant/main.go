package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/SultanKs4/wassistant/config"
	"github.com/SultanKs4/wassistant/internal/whatsapp"
	"github.com/SultanKs4/wassistant/pkg/db/postgresql"
	"github.com/SultanKs4/wassistant/pkg/db/redis"
)

func main() {
	cfg, err := config.NewConfig("../../config/config-local.yml")
	if err != nil {
		panic(err)
	}
	pg, err := postgresql.NewPg(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	if err := pg.Migrate(); err != nil {
		panic(err)
	}
	defer pg.Disconnect()

	redb, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer redb.Close()

	// aws, err := aws.NewAws(cfg.Aws)
	// if err != nil {
	// 	panic(err)
	// }

	wa, err := whatsapp.Connect(pg)
	if err != nil {
		panic(err)
	}
	defer wa.Disconnect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
