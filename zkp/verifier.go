package zkp

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp/algorithm"
	"github.com/mindaugasrukas/zkp_example/zkp/pedersen"
)

var (
	ChallengeError = errors.New("challenge construction error")
)

type (
	PedersenVerifier struct {
		pedersen.Verifier

		challenge *big.Int
	}
)

func NewVerifier() *PedersenVerifier {
	return &PedersenVerifier{
		Verifier: pedersen.Verifier{
			P: big.NewInt(P),
			Q: big.NewInt(Q),
			G: big.NewInt(G),
			H: big.NewInt(H),
		},
	}
}

// CreateAuthenticationChallenge Creates a challenge that the client has to answer
func (v *PedersenVerifier) CreateAuthenticationChallenge() (challenge *big.Int, err error) {
	challenge, err = rand.Int(rand.Reader, v.Q)
	if err != nil {
		return nil, ChallengeError
	}
	if challenge == algorithm.ZERO {
		// if zero repeat
		return v.CreateAuthenticationChallenge()
	}
	v.challenge = challenge
	return challenge, nil
}

// VerifyAuthentication verifies the answer received from the client against the commits
// Verify the answer s using this formula:
// g and h is public keys
// c is the challenge
// y1, y2 is client registered commits
// r1, r2 is client authentication commits
// r1 = g^s * y1^c  AND  r2 = h^s * y2^c
func (v *PedersenVerifier) VerifyAuthentication(commits *Commits, authRequest *Commits, answer *big.Int) bool {
	var y1, y2 big.Int

	p := big.NewInt(P)
	g := big.NewInt(G)
	h := big.NewInt(H)
	log.Printf("y1=%v, y2=%v, g=%v, h=%v, answer=%v, challenge=%v", commits.C1, commits.C2, g, h, answer, v.challenge)

	g.Exp(g, answer, nil)
	y1.Exp(commits.C1, v.challenge, nil)
	log.Printf("g^answer=%v, y1^challenge=%v", g, &y1)
	var result1 big.Int
	result1.Mul(g, &y1)
	result1.Mod(&result1, p)
	log.Print("result1: (g^answer)*(y1^challenge) = ", &result1)

	h.Exp(h, answer, nil)
	y2.Exp(commits.C2, v.challenge, nil)
	log.Printf("h^answer=%v, y2^challenge=%v", h, &y2)
	var result2 big.Int
	result2.Mul(h, &y2)
	result2.Mod(&result2, p)
	log.Print("result2: (h^answer)*(y2^challenge) = ", &result2)

	log.Printf("r1=%v, r2=%v", authRequest.C1, authRequest.C2)
	log.Printf("(result1==r1)=%v, (result2==r2)=%v", result1.Cmp(authRequest.C1) == 0, result2.Cmp(authRequest.C2) == 0)

	return result1.Cmp(authRequest.C1) == 0 && result2.Cmp(authRequest.C2) == 0
}
