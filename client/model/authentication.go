package model

import (
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

// GetChallenge translates request to internal type
func GetChallenge(challengeResponse *zkp_pb.ChallengeResponse) *big.Int {
	var challenge big.Int
	challenge.SetBytes(challengeResponse.GetChallenge())
	return &challenge
}
