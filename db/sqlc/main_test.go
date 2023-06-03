package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable"
)

func TestMain(m *testing.M) {

	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
