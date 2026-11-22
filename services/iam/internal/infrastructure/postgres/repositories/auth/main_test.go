package auth

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"todolist/internal/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	Postgres *pgxpool.Pool
)

func loadTestEnv() {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"..",
		"..",
		"..",
	)
	envFileName := ".env.test"
	err := godotenv.Load(
		filepath.Join(root, envFileName),
	)
	if err != nil {
		log.Fatalf("load test %v: %v", envFileName, err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	loadTestEnv()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(os.Getenv("DATABASE_NAME")),
		postgres.WithUsername(os.Getenv("DATABASE_USER")),
		postgres.WithPassword(os.Getenv("DATABASE_PASSWORD")),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	connString, err := postgresContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	if err != nil {
		log.Fatalf("get connection string: %v", err)
	}

	// Run migrations once
	err = runMigrations(connString)
	if err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	// Create pool used by tests
	Postgres, err = pgxpool.New(
		ctx,
		connString,
	)
	if err != nil {
		log.Fatalf("create pool: %v", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func runMigrations(url string) error {
	source, err := iofs.New(
		database.MigrationFiles,
		"migrations",
	)

	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		url,
	)

	if err != nil {
		return err
	}

	return m.Up()
}
