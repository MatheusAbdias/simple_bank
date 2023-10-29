package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/MatheusAbdias/go_simple_bank/api"
	"github.com/MatheusAbdias/go_simple_bank/cmd"
	"github.com/MatheusAbdias/go_simple_bank/config"
	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	conn, err := sql.Open(config.Driver, config.Source)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	cmd.Migrate()

	store := db.NewSQLStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}
