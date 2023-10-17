package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/MatheusAbdias/go_simple_bank/util"
)

var testQueries *Queries

var dbConn *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, container := util.SetupDataBase(ctx)
	dbConn = db
	testQueries = New(db)

	exitCode := m.Run()

	util.TearDownDataBase(ctx, db, container)

	os.Exit(exitCode)
}
