package util

import (
	"context"
	"testing"
)

func TestConnectWithDB(t *testing.T) {
	ctx := context.Background()

	container, err := NewPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("error creating container: %s", err)
	}
	defer container.Terminate(ctx)

	db, err := NewDB(container.ConnectionString)
	if err != nil {
		t.Fatalf("error creating database: %s", err)
	}
	defer db.Close()

}

func TestHealthChecker(t *testing.T) {
	ctx := context.Background()

	container, err := NewPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("error creating container: %s", err)
	}
	defer container.Terminate(ctx)

	db, err := NewDB(container.ConnectionString)
	if err != nil {
		t.Fatalf("error creating database: %s", err)
	}
	defer db.Close()

	healthChecker := NewHealthChecker(db)
	err = healthChecker.CheckHealth(ctx)
	if err != nil {
		t.Fatalf("error checking health: %s", err)
	}
}
