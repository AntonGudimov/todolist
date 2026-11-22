package postgres

import (
	"context"
	"todolist/internal/domain"
	"todolist/internal/infrastructure/postgres/generated"
	"todolist/internal/infrastructure/postgres/repositories/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUoW struct {
	db *pgxpool.Pool
	q  *generated.Queries
}

func NewPostgres(db *pgxpool.Pool) domain.UnitOfWork {
	return &PostgresUoW{
		db: db,
		q:  generated.New(db),
	}
}

func (u PostgresUoW) Auth() domain.AuthRepository {
	return auth.NewAuthPostgresRepository(u.q)
}

func (u *PostgresUoW) Do(ctx context.Context, fn func(domain.UnitOfWork) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	txUow := &PostgresUoW{
		db: u.db,
		q:  generated.New(tx),
	}

	if err := fn(txUow); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
