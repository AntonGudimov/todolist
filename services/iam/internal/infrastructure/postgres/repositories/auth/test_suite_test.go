package auth

import (
	"context"
	"testing"
	"todolist/internal/domain"
	"todolist/internal/infrastructure/postgres/generated"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type AuthPostgresRepositorySuite struct {
	suite.Suite

	tx      pgx.Tx
	queries *generated.Queries
	repo    domain.AuthRepository

	ctx context.Context
}

func (s *AuthPostgresRepositorySuite) SetupTest() {
	s.ctx = context.Background()
	tx, err := Postgres.Begin(s.ctx)
	s.Require().NoError(err)

	s.tx = tx
	s.queries = generated.New(tx)
	s.repo = NewAuthPostgresRepository(s.queries)

}

func (s *AuthPostgresRepositorySuite) TearDownTest() {
	_ = s.tx.Rollback(s.ctx)
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthPostgresRepositorySuite))
}
