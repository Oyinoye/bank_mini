package main

import (
	"database/sql"
	"log"

    "github.com/Oyinoye/bank_mini/util"

	"github.com/Oyinoye/bank_mini/api"
	db "github.com/Oyinoye/bank_mini/db/sqlc"
	_ "github.com/lib/pq"
)


func main() {

    config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	// if config.Environment == "development" {
	// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// }


	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
    server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot start server:", err)
    }
}
