package zkp_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/stretchr/testify/assert"
)

func TestCreateRegisterCommits(t *testing.T) {
	assert := assert.New(t)
	prover := zkp.NewProver(123)
	commits, err := prover.CreateRegisterCommits()
	assert.NoError(err)
	assert.Equal(int64(16), commits.C1.Int64())
	assert.Equal(int64(12), commits.C2.Int64())
}


func TestCreateAuthenticationCommits(t *testing.T) {
	assert := assert.New(t)
	prover := zkp.NewProver(123)
	commits, err := prover.CreateAuthenticationCommits()
	assert.NoError(err)
	assert.True(commits.C1.Int64() > 0)
	assert.True(commits.C2.Int64() > 0)
}

func TestProveAuthentication(t *testing.T) {
	assert := assert.New(t)
	prover := zkp.NewProver(123)
	_, err := prover.CreateAuthenticationCommits()
	assert.NoError(err)
	answer := prover.ProveAuthentication(big.NewInt(5))
	assert.True(answer.Int64() > 0)
}
