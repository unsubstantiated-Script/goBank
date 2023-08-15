package main

import (
	"database/sql"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"goBank/api"
	db "goBank/db/sqlc"
	"goBank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Can't connect to DB!")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start server!", err)
	}
}
