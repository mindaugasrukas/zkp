package zkp_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthenticationChallenge(t *testing.T) {
	assert := assert.New(t)
	verifier := zkp.NewVerifier()
	challenge, err := verifier.CreateAuthenticationChallenge()
	assert.NoError(err)
	assert.True(challenge.Int64() > 0)
}

func TestVerifyAuthentication_Success(t *testing.T) {
	assert := assert.New(t)
	verifier := zkp.NewVerifier()
	commits := &zkp.Commits{
		C1: big.NewInt(2),
		C2: big.NewInt(3),
	}
	authRequest := &zkp.Commits{
		C1: big.NewInt(8),
		C2: big.NewInt(4),
	}
	challenge := big.NewInt(4)
	answer := big.NewInt(5)
	result := verifier.VerifyAuthentication(commits, authRequest, challenge, answer)
	assert.True(result)
}

func TestVerifyAuthentication_Fail(t *testing.T) {
	assert := assert.New(t)
	verifier := zkp.NewVerifier()
	commits := &zkp.Commits{
		C1: big.NewInt(2),
		C2: big.NewInt(3),
	}
	authRequest := &zkp.Commits{
		C1: big.NewInt(8),
		C2: big.NewInt(4),
	}
	challenge := big.NewInt(4)
	// report the wrong answer
	answer := big.NewInt(55)
	result := verifier.VerifyAuthentication(commits, authRequest, challenge, answer)
	assert.False(result)
}
