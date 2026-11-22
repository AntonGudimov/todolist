package fakes

import "todolist/internal/domain"

type FakeKeyStore struct {
	key             domain.Key
	getActiveKeyErr error
}

func NewFakeKeyStore(key domain.Key, getActiveKeyErr error) *FakeKeyStore {
	return &FakeKeyStore{key: key, getActiveKeyErr: getActiveKeyErr}
}

func (ks *FakeKeyStore) GetSigningKey() (domain.Key, error) {
	if ks.getActiveKeyErr != nil {
		return domain.Key{}, ks.getActiveKeyErr
	}
	return ks.key, nil
}

func (ks *FakeKeyStore) Reload() error {
	return nil
}
