package fakes

import (
	"todolist/internal/domain"

	"github.com/google/uuid"
)

type FakeDB struct {
	Users               map[uuid.UUID]*domain.User
	EmailVerifications  map[uuid.UUID]*domain.EmailVerification
	RefreshTokens       map[uuid.UUID]*domain.RefreshToken
	PasswordResetTokens map[uuid.UUID]*domain.PasswordResetToken
}

func NewFakeDB() *FakeDB {
	return &FakeDB{
		Users:               make(map[uuid.UUID]*domain.User, 0),
		EmailVerifications:  make(map[uuid.UUID]*domain.EmailVerification, 0),
		RefreshTokens:       make(map[uuid.UUID]*domain.RefreshToken, 0),
		PasswordResetTokens: make(map[uuid.UUID]*domain.PasswordResetToken, 0),
	}
}

func (db *FakeDB) clearUsers() {
	for u := range db.Users {
		delete(db.Users, u)
	}
}

func (db *FakeDB) clearEmailVerifications() {
	for e := range db.EmailVerifications {
		delete(db.EmailVerifications, e)
	}
}

func (db *FakeDB) clearRefreshTokens() {
	for r := range db.RefreshTokens {
		delete(db.RefreshTokens, r)
	}
}

func (db *FakeDB) clearPasswordResetTokens() {
	for p := range db.PasswordResetTokens {
		delete(db.PasswordResetTokens, p)
	}
}

func (db *FakeDB) Clear() {
	db.clearUsers()
	db.clearEmailVerifications()
	db.clearRefreshTokens()
	db.clearPasswordResetTokens()
}
