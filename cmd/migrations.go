package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const ProjectName = "simple_bank"

func Migrate() {
	dbUrl := os.Getenv("DATABASE_URL")
	baseDir := os.Getenv("BASE_DIR")

	m, err := migrate.New(fmt.Sprintf("file://%s/db/migration", baseDir), dbUrl)

	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
