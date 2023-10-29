package util

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/MatheusAbdias/go_simple_bank/cmd"
)

func SetupDataBase(ctx context.Context) (*sql.DB, *PostgresContainer) {
	container, err := NewPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("error creating container: %s", err)
	}

	db, err := NewDB(container.ConnectionString)
	if err != nil {
		log.Fatalf("error creating database: %s", err)
	}

	healthChecker := NewHealthChecker(db)
	err = healthChecker.CheckHealth(ctx)
	if err != nil {
		log.Fatalf("error checking health: %s", err)
	}
	dir := cmd.FindProjectDir()
	os.Setenv("BASE_DIR", dir)

	os.Setenv("DATABASE_URL", container.ConnectionString)
	cmd.Migrate()

	return db, container
}

func TearDownDataBase(ctx context.Context, db *sql.DB, container *PostgresContainer) {
	db.Close()
	container.Terminate(ctx)
}
