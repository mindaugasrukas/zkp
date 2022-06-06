package zkp

import "math/big"

// Public information
const (
	P = 23
	G = 4
	H = 9
	Q = 11
)

type (
	// Commits is user commitments for registration or authentication
	Commits struct {
		C1, C2 *big.Int
	}

	// UUID is Unique User ID
	UUID string
)
