package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"todolist/internal/database"
)

func main() {
	databaseURL := os.Getenv("IAM_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("IAM_DATABASE_URL is required")
	}

	source, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance(
		"embedded",
		source,
		databaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
