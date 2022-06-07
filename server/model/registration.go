package model

import (
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

// GetRegistration translates request commits to internal types
func GetRegistration(registerRequest *zkp_pb.RegisterRequest) (user zkp.UUID, commits *zkp.Commits) {
	var y1, y2 big.Int
	c := registerRequest.GetCommits()[0]

	y1.SetBytes(c.GetY1())
	y2.SetBytes(c.GetY2())

	user = zkp.UUID(registerRequest.GetUser())
	return user, &zkp.Commits{
		C1: &y1,
		C2: &y2,
	}
}
