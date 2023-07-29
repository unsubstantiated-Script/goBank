package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"goBank/util"
	"log"
	"os"
	"testing"
)

/*
ATOMICITY: All operations complete successfully or transaction fails and db is unchanged
Consistency: db state is valid after transaction. All constraints satisfied.
Isolation: Concurrent transactions do not affect each other.
Durability: Data written by a successful transaction must be recorded in persistent storage.
*/

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("can't load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Can't connect to DB!")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
