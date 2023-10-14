package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/MatheusAbdias/go_simple_bank/cmd"
	"github.com/MatheusAbdias/go_simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var dbConn *sql.DB

func findProjectDirRecursive(currentDir, projectName string) string {
	if base := filepath.Base(currentDir); base == projectName {
		return currentDir
	}

	if currentDir == "/" {
		return ""
	}

	parentDir := filepath.Dir(currentDir)
	return findProjectDirRecursive(parentDir, projectName)
}
func findProjectDir(projectName string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return findProjectDirRecursive(currentDir, projectName)
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, container := setupDataBase(ctx)
	dbConn = db
	testQueries = New(db)

	exitCode := m.Run()

	tearDownDataBase(ctx, db, container)

	os.Exit(exitCode)
}

func setupDataBase(ctx context.Context) (*sql.DB, *util.PostgresContainer) {
	container, err := util.NewPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("error creating container: %s", err)
	}

	db, err := util.NewDB(container.ConnectionString)
	if err != nil {
		log.Fatalf("error creating database: %s", err)
	}

	healthChecker := util.NewHealthChecker(db)
	err = healthChecker.CheckHealth(ctx)
	if err != nil {
		log.Fatalf("error checking health: %s", err)
	}
	dir := findProjectDir(cmd.ProjectName)
	os.Setenv("BASE_DIR", dir)

	os.Setenv("DATABASE_URL", container.ConnectionString)
	cmd.Migrate()

	return db, container
}

func tearDownDataBase(ctx context.Context, db *sql.DB, container *util.PostgresContainer) {
	db.Close()
	container.Terminate(ctx)
}
