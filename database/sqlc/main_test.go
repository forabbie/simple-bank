package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/forabbie/vank-app/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Ensure New(conn) is defined in your package and initializes Queries
	testQueries = New(testDB)

	os.Exit(m.Run())
}
