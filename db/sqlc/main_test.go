package db

import (
	"database/sql"
	_ "github.com/lib/pq"
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
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Can't connect to DB!")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
