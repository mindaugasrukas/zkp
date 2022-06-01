package zkp

import "math/big"

// Public information
const (
	G = 3
	H = 11
	Q = 33
)

type (
	// Commits is user secret commits
	Commits struct {
		Y1, Y2 *big.Int
	}

	// AuthenticationRequest is user auth request
	AuthenticationRequest struct {
		R1, R2 *big.Int
	}

	// Challenge is server challenge
	Challenge *big.Int

	// Answer is user answer
	Answer *big.Int

	// UUID is Unique User ID
	UUID string
)
