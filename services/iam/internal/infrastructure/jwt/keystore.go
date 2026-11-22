package jwt

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"todolist/internal/domain"
)

type FileKeyStore struct {
	mu sync.RWMutex

	dir          string
	signingKeyID string
	keys         map[string]domain.Key
}

func NewFileKeyStore(dir string) (*FileKeyStore, error) {
	ks := &FileKeyStore{
		dir:  dir,
		keys: make(map[string]domain.Key, 0),
	}
	if err := ks.Reload(); err != nil {
		return nil, err
	}
	return ks, nil
}

func (ks *FileKeyStore) Reload() error {
	activeKeyID, activeKeyIDErr := loadActiveKeyID(ks.dir)
	if activeKeyIDErr != nil {
		return activeKeyIDErr
	}
	validateActiveKeyIDErr := validateActiveKeyID(activeKeyID, ks.dir)
	if validateActiveKeyIDErr != nil {
		return validateActiveKeyIDErr
	}
	entries, entriesErr := os.ReadDir(ks.dir)
	if entriesErr != nil {
		return fmt.Errorf("read dir: %w", entriesErr)
	}
	keys := make(map[string]domain.Key)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		keyID := entry.Name()
		key, keyErr := loadKeyFromDir(filepath.Join(ks.dir, keyID))
		if keyErr != nil {
			return keyErr
		}
		keys[keyID] = key
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.signingKeyID = activeKeyID
	ks.keys = keys
	return nil
}

func (ks *FileKeyStore) GetSigningKey() (domain.Key, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	key := ks.keys[ks.signingKeyID]
	return key, nil
}
