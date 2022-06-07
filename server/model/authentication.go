package model

import (
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

// GetAuthentication translates request commits to internal types
func GetAuthentication(authRequest *zkp_pb.AuthRequest) (user zkp.UUID, commits *zkp.Commits) {
	var r1, r2 big.Int
	c := authRequest.GetCommits()[0]

	r1.SetBytes(c.GetR1())
	r2.SetBytes(c.GetR2())

	user = zkp.UUID(authRequest.GetUser())
	return user, &zkp.Commits{
		C1: &r1,
		C2: &r2,
	}
}

// GetAnswer translates request to internal type
func GetAnswer(answerRequest *zkp_pb.AnswerRequest) *big.Int {
	var answer big.Int
	answer.SetBytes(answerRequest.GetAnswer())
	return &answer
}
