package model_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/server/model"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthentication(t *testing.T) {
	assert := assert.New(t)
	authRequest := &zkp_pb.AuthRequest{
		User: "test-user",
		Commits: []*zkp_pb.AuthRequest_Commits{
			{
				R1: []byte{0xc},
				R2: []byte{0xd},
			},
		},
	}
	user, commits := model.GetAuthentication(authRequest)
	assert.Equal(zkp.UUID("test-user"), user)
	expectedCommits := &zkp.Commits{
		C1: big.NewInt(12),
		C2: big.NewInt(13),
	}
	assert.Equal(expectedCommits, commits)
}

func TestGetAnswer(t *testing.T) {
	assert := assert.New(t)
	answerRequest := &zkp_pb.AnswerRequest{
		Answer: []byte{0xd},
	}
	answer := model.GetAnswer(answerRequest)
	assert.Equal(big.NewInt(13), answer)
}
