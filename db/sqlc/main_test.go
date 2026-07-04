package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Oyinoye/bank_mini/util"
	// _ "github.com/lib/pq"
	"github.com/jackc/pgx/v5/pgxpool"
)

// const (
//     dbDriver = "postgres"
//     dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// )

// var testStore *Queries
// var testDB *sql.DB

// func TestMain(m *testing.M) {
//     config, err := util.LoadConfig("../..")
//     if err != nil {
//         log.Fatal("cannot load convig:", err)
//     }

// 	testDB, err = sql.Open(config.DBDriver, config.DBSource)
// 	if err != nil {
// 		log.Fatal("cannot connect to db:", err)
// 	}

// 	testStore = New(testDB)

// 	os.Exit(m.Run())
// }

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
