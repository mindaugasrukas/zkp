package zkp

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ChallengeError = errors.New("challenge construction error")
)

type Verifier struct {
	challenge Challenge
}

// CreateAuthenticationChallenge Creates a challenge that the client has to answer
func (v *Verifier) CreateAuthenticationChallenge(authRequest AuthenticationRequest) (Challenge, error) {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		return nil, ChallengeError
	}
	r := big.NewInt(int64(b[0]))
	v.challenge = r
	return r, nil
}

// VerifyAuthentication verifies the answer received from the client against the commits
func (v *Verifier) VerifyAuthentication(commits Commits, authRequest AuthenticationRequest, answer Answer) bool {
	y1 := commits.Y1
	y2 := commits.Y2

	g := big.NewInt(G)
	h := big.NewInt(H)

	g.Exp(g, answer, nil)
	y1.Exp(y1, v.challenge, nil)
	var result1 big.Int
	result1.Mul(g, y1)

	h.Exp(h, answer, nil)
	y2.Exp(y2, v.challenge, nil)
	var result2 big.Int
	result2.Mul(g, y1)

	return result1.Cmp(authRequest.R1) == 0 && result2.Cmp(authRequest.R2) == 0
}
