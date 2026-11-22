package fakes

import (
	"context"
	"todolist/internal/domain"
)

type FakeUnitOfWork struct {
	Db *FakeDB
}

func NewFakeUnitOfWork(db *FakeDB) *FakeUnitOfWork {
	return &FakeUnitOfWork{Db: db}
}

func (u *FakeUnitOfWork) Auth() domain.AuthRepository {
	return NewFakeAuthRepository(u.Db)
}

func (u *FakeUnitOfWork) Do(ctx context.Context, fn func(domain.UnitOfWork) error) error {
	return fn(u)
}
