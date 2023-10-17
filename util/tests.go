package util

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/MatheusAbdias/go_simple_bank/cmd"
)

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
	dir := findProjectDir(cmd.ProjectName)
	os.Setenv("BASE_DIR", dir)

	os.Setenv("DATABASE_URL", container.ConnectionString)
	cmd.Migrate()

	return db, container
}

func TearDownDataBase(ctx context.Context, db *sql.DB, container *PostgresContainer) {
	db.Close()
	container.Terminate(ctx)
}
