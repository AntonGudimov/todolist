package domain

import (
	"context"
)

type UnitOfWork interface {
	Auth() AuthRepository
	Do(ctx context.Context, fn func(UnitOfWork) error) error
}
