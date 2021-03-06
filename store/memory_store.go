package store

import (
	"errors"

	"github.com/mindaugasrukas/zkp_example/zkp"
)

var (
	UserExistsError       = errors.New("user already exists")
	UserDoesNotExistError = errors.New("user doesn't exist")
)

type (
	InMemoryStore struct {
		store map[zkp.UUID]*zkp.Commits
	}
)

// NewInMemoryStore returns a new store instance
func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{
		store: make(map[zkp.UUID]*zkp.Commits),
	}
}

// Add user to the store
// returns UserExistsError if user already exists
func (m InMemoryStore) Add(user zkp.UUID, commits *zkp.Commits) error {
	if _, ok := m.store[user]; ok {
		return UserExistsError
	}
	m.store[user] = commits
	return nil
}

// Get user data from the store
// returns UserDoesNotExistError if user doesn't exist
func (m InMemoryStore) Get(user zkp.UUID) (*zkp.Commits, error) {
	data, ok := m.store[user]
	if !ok {
		return nil, UserDoesNotExistError
	}
	return data, nil
}
