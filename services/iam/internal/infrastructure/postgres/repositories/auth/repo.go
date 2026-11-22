package auth

import (
	"todolist/internal/domain"
	"todolist/internal/infrastructure/postgres/generated"
)

type AuthPostgresRepository struct {
	q *generated.Queries
}

func NewAuthPostgresRepository(q *generated.Queries) domain.AuthRepository {
	return &AuthPostgresRepository{q: q}
}
