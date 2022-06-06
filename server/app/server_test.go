package app_test

import (
	"math/big"
	"testing"

	svr "github.com/mindaugasrukas/zkp_example/server/app"
	"github.com/mindaugasrukas/zkp_example/store"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/stretchr/testify/assert"
)

func TestServer_Register(t *testing.T) {
	assert := assert.New(t)
	server := svr.NewServer()
	user := zkp.UUID("userid-123")
	commits := zkp.Commits{
		C1: big.NewInt(int64(123)),
		C2: big.NewInt(int64(345)),
	}

	// register a new user without the error
	err := server.Register(user, &commits)
	assert.NoError(err)

	// fail to register duplicate user
	err = server.Register(user, &commits)
	assert.ErrorIs(err, store.UserExistsError)
}

func TestServer_CreateAuthenticationChallenge(t *testing.T) {
	assert := assert.New(t)
	server := svr.NewServer()
	user := zkp.UUID("userid-123")
	commits := zkp.Commits{
		C1: big.NewInt(int64(123)),
		C2: big.NewInt(int64(345)),
	}

	// fail to initiate auth session
	_, err := server.Verifier.CreateAuthenticationChallenge()
	assert.NoError(err)

	// register dummy user
	err = server.Register(user, &commits)
	assert.NoError(err)

	// return non empty auth challenge
	challenge, err := server.Verifier.CreateAuthenticationChallenge()
	assert.NoError(err)
	assert.True(challenge != nil)
}
