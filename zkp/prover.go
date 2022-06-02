package zkp

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Prover struct {
	k *big.Int
	password int64
}

var (
	AuthenticationRequestError = errors.New("error creating authentication request")
)

func NewProver(password int64) *Prover {
	return &Prover{
		password: password,
	}
}


// CreateRegisterCommits Creates the commits to register the user in the server
// Having secret x and public keys g and h,
// calculate y1 = g^x and y2 = h^x
func (p *Prover) CreateRegisterCommits() Commits {
	var y1, y2 big.Int
	y1.Exp(big.NewInt(G), big.NewInt(p.password), nil)
	y2.Exp(big.NewInt(H), big.NewInt(p.password), nil)

	return Commits{
		Y1: &y1,
		Y2: &y2,
	}
}

// CreateAuthenticationRequest Creates an authentication request to start the authentication against the Server
// Generate random k and using public keys g and h,
// calculate r1 = g^k and r2 = h^k
func (p *Prover) CreateAuthenticationRequest() (*AuthenticationRequest, error) {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		return nil, AuthenticationRequestError
	}
	p.k = big.NewInt(int64(b[0]))

	var r1, r2 big.Int
	r1.Exp(big.NewInt(G), p.k, nil)
	r2.Exp(big.NewInt(H), p.k, nil)

	return &AuthenticationRequest{
		R1: &r1,
		R2: &r2,
	}, nil
}

// ProveAuthentication Returns the answer to the challenge
// Having secret x and random k
// given challenge c and using public q
// calculate the answer s = k - c * x (mod q)
func (p *Prover) ProveAuthentication(challenge Challenge) Answer {
	var c, answer big.Int
	c.Mul(challenge, big.NewInt(p.password))
	answer.Sub(p.k, &c)
	answer.Mod(&answer, big.NewInt(Q))
	return &answer
}
