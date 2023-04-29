package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/SultanKs4/wassistant/internal/whatsapp"
	"github.com/SultanKs4/wassistant/pkg/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	pg, err := db.NewPg(os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		panic(err)
	}
	pg.Migrate()
	defer pg.Disconnect()

	dbInt, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	redb, err := db.NewRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), dbInt)
	if err != nil {
		panic(err)
	}
	defer redb.Close()

	db, err := pg.DB.DB()
	if err != nil {
		panic(err)
	}
	wa, err := whatsapp.Connect(db)
	if err != nil {
		panic(err)
	}
	defer wa.Disconnect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
