package zkp

import (
	"log"
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp/algorithm"
	"github.com/mindaugasrukas/zkp_example/zkp/pedersen"
)

type (
	PedersenProver struct {
		pedersen.Prover
	}
)

func NewProver(password int64) *PedersenProver {
	q := big.NewInt(Q)
	private := &algorithm.Zr{
		Value:  big.NewInt(password),
		Modulo: q,
	}
	return &PedersenProver{
		Prover: pedersen.Prover{
			P: big.NewInt(P),
			Q: q,
			G: big.NewInt(G),
			H: big.NewInt(H),
			X: private,
			R: private,
		},
	}
}

// CreateRegisterCommits Creates the commits to register the user in the server
// Having secret x and public keys g and h,
// calculate y1 = g^x and y2 = h^x
func (p *PedersenProver) CreateRegisterCommits() (*Commits, error) {
	y1 := big.NewInt(0)
	y1.Exp(p.G, p.X.Value, p.P)
	log.Print("g = ", p.G)
	log.Print("g^rx = ", y1)

	y2 := big.NewInt(0)
	y2.Exp(p.H, p.X.Value, p.P)
	log.Print("h = ", p.H)
	log.Print("h^rr = ", y2)

	return &Commits{
		C1: y1,
		C2: y2,
	}, nil
}

// CreateAuthenticationCommits Creates an authentication request to start the authentication against the Server
// Generate random k and using public keys g and h,
// calculate r1 = g^k and r2 = h^k
func (p *PedersenProver) CreateAuthenticationCommits() (*Commits, error) {
	r1, r2, err := p.Commits()
	if err != nil {
		return nil, err
	}

	return &Commits{
		C1: r1,
		C2: r2,
	}, nil
}

// ProveAuthentication Returns the answer to the challenge
// Having secret x and random k
// given challenge c and using public q
// calculate the answer s = k - c * x (mod q)
func (p *PedersenProver) ProveAuthentication(challenge *big.Int) (answer *big.Int) {
	return p.Prove(challenge)[0]
}
