package util

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	return db, nil
}
