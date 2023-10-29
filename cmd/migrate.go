package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const ProjectName = "simple_bank"

func Migrate() {
	dbUrl := os.Getenv("DATABASE_URL")
	baseDir := os.Getenv("BASE_DIR")

	fullPath := fmt.Sprintf("file://%s/db/migration", baseDir)
	m, err := migrate.New(fullPath, dbUrl)

	if err != nil {
		log.Fatal(fullPath)
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func findProjectDirRecursive(currentDir string) string {
	if base := filepath.Base(currentDir); base == ProjectName {
		return currentDir
	}

	if currentDir == "/" {
		return ""
	}

	parentDir := filepath.Dir(currentDir)
	return findProjectDirRecursive(parentDir)
}

func FindProjectDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return findProjectDirRecursive(currentDir)
}
