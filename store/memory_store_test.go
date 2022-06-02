package store_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/store"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore_Add(t *testing.T) {
	assert := assert.New(t)
	registry := store.NewInMemoryStore()

	// Add a new user without the error
	user := zkp.UUID("userid-123")
	err := registry.Add(user, &zkp.Commits{
		Y1: big.NewInt(int64(123)),
		Y2: big.NewInt(int64(345)),
	})
	assert.NoError(err)

	// Fail to add a duplicate user
	err = registry.Add(user, &zkp.Commits{
		Y1: big.NewInt(int64(789)),
		Y2: big.NewInt(int64(567)),
	})
	assert.ErrorIs(err, store.UserExistsError)

	// Add a new user without the error
	user2 := zkp.UUID("userid-789")
	err = registry.Add(user2, &zkp.Commits{
		Y1: big.NewInt(int64(123)),
		Y2: big.NewInt(int64(345)),
	})
	assert.NoError(err)
}

func TestInMemoryStore_Get(t *testing.T) {
	assert := assert.New(t)
	registry := store.NewInMemoryStore()

	// Fail to get non existing user
	user := zkp.UUID("userid-123")
	_, err := registry.Get(user)
	assert.ErrorIs(err, store.UserDoesNotExistError)

	// Add a dummy user
	commits := &zkp.Commits{
		Y1: big.NewInt(int64(123)),
		Y2: big.NewInt(int64(345)),
	}
	err = registry.Add(user, commits)
	assert.NoError(err)

	// Successfully get existing user data
	data, err := registry.Get(user)
	assert.NoError(err)
	assert.Equal(commits, data)
}
