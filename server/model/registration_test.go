package model_test

import (
	"math/big"
	"testing"

	"github.com/mindaugasrukas/zkp_example/server/model"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistration(t *testing.T) {
	assert := assert.New(t)
	registerRequest := &zkp_pb.RegisterRequest{
		User: "test-user",
		Commits: []*zkp_pb.RegisterRequest_Commits{
			{
				Y1: []byte{0xc},
				Y2: []byte{0xd},
			},
		},
	}
	user, commits := model.GetRegistration(registerRequest)
	assert.Equal(zkp.UUID("test-user"), user)
	expectedCommits := &zkp.Commits{
		C1: big.NewInt(12),
		C2: big.NewInt(13),
	}
	assert.Equal(expectedCommits, commits)
}
