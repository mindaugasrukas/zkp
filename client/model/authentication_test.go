package model_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/client/model"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/stretchr/testify/assert"
)

func TestGetChallenge(t *testing.T) {
	assert := assert.New(t)
	challengeResponse := &zkp_pb.ChallengeResponse{
		Challenge: []byte{0xd},
	}
	challenge := model.GetChallenge(challengeResponse)
	assert.Equal(big.NewInt(13), challenge)
}
