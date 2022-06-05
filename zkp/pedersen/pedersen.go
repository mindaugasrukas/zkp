package pedersen

//
// Copied from https://github.com/lziest/goZKP
//

import (
	"crypto/sha1"
	"log"
	"math/big"

	"github.com/mindaugasrukas/zkp_example/zkp/algorithm"
)

// Prover proves the knowledge of (x, r) such that z = (g**x) * (h**r).
type Prover struct {
	P *big.Int  // Zp as Group
	Q *big.Int  // G's order
	G *big.Int  // Group generator g
	H *big.Int  // Group generator h
	X *algorithm.Zr // Private x
	R *algorithm.Zr // Private r
}

type Verifier struct {
	P *big.Int // Zp as Group
	Q *big.Int // G's order
	G *big.Int // Group generator g
	H *big.Int // Group generator h
	// public commitments
	Y1 *big.Int	// y1 = (g**x)
	Y2 *big.Int // y2 = (h**r)
	Z *big.Int  // z = (g**x) * (h**r) = y1 * y2
}

func Hash(q *big.Int, values ...*big.Int) *big.Int {
	hash := sha1.New()
	for _, value := range values {
		hash.Write(value.Bytes())
	}
	h := hash.Sum(nil)
	ret := big.NewInt(0)
	ret.SetBytes(h)
	ret.Mod(ret, q)
	return ret
}

func (p *Prover) Commits() (y1 *big.Int, y2 *big.Int, err error) {
	rx, err := p.X.Commit()
	if err != nil {
		return nil, nil, err
	}
	rr, err := p.R.Commit()
	if err != nil {
		return nil, nil, err
	}

	g := big.NewInt(0)
	g.Exp(p.G, rx, p.P)
	log.Print("g = ", p.G.Int64())
	log.Print("g^rx = ", g.Int64())

	h := big.NewInt(0)
	h.Exp(p.H, rr, p.P)
	log.Print("h = ", p.H.Int64())
	log.Print("h^rr = ", h.Int64())

	return g, h, nil
}

func (p *Prover) Commit() ([]*big.Int, error) {
	y1, y2, err := p.Commits()
	if err != nil {
		return nil, err
	}
	comm := big.NewInt(0)
	comm.Mul(y1, y2)
	comm.Mod(comm, p.P)
	return []*big.Int{comm}, nil
}

func (p *Prover) Prove(c *big.Int) []*big.Int {
	if c.Cmp(p.Q) > 0 {
		c.Mod(c, p.Q)
	}

	// tx = (rx - c * p.X) mod Q
	// tr = (rr - c * p.X) mod Q
	tx := p.X.Prove(c)
	tr := p.R.Prove(c)
	ret := []*big.Int{tx, tr}
	return ret
}

func (p *Prover) Sign(m *big.Int) ([]*big.Int, error) {
	comm, err := p.Commit()
	if err != nil {
		return nil, err
	}
	comm = append([]*big.Int{m}, comm...)
	c := Hash(p.Q, comm...)
	proof := p.Prove(c)
	proof = append(proof, c)
	return proof, nil
}

// ConsumeVerify To enable aggregated response
// verification, i.e. multiple provers proving the knowledge of
// their commitments at the same time via one challenge, ConsumeVerify
// will only consume comm and resp as a stream and will return
// remaining commitments and remaining response as remComm and
// remResp respectively.
func (v *Verifier) ConsumeVerify(comm []*big.Int, c *big.Int, resp []*big.Int) (valid bool, remComm, remResp []*big.Int) {
	valid = false
	if len(resp) < 2 {
		return false, comm, resp
	}

	if len(comm) < 1 {
		return false, comm, resp
	}

	var rv *big.Int
	rv, remResp = v.RecoverCommitment(c, resp)

	rc := comm[0]
	remComm = comm[1:]

	if rv.Cmp(rc) != 0 {
		return
	}
	valid = true
	return
}

// Verify would verify commitments against the challenge value
// and prover's response values.
func (v *Verifier) Verify(comm []*big.Int, c *big.Int, resp []*big.Int) bool {
	valid, remComm, remResp := v.ConsumeVerify(comm, c, resp)

	if len(remComm) != 0 || len(remResp) != 0 || !valid {
		return false
	}

	return true
}

func (v *Verifier) RecoverCommitment(c *big.Int, resp []*big.Int) (rc *big.Int, remResp []*big.Int) {
	sx := resp[0]
	sr := resp[1]
	remResp = resp[2:]

	// rc = G**sx * H**sr * Z**c
	//    = G**(sx + x * c) * H**(sr + r * c)
	//    = G**rx * H**rr (mod P)
	rc = big.NewInt(0)
	rc.Exp(v.G, sx, v.P)

	tmp := big.NewInt(0)
	tmp.Exp(v.H, sr, v.P)
	rc.Mul(rc, tmp)
	rc.Mod(rc, v.P)

	tmp.Exp(v.Z, c, v.P)
	rc.Mul(rc, tmp)
	rc.Mod(rc, v.P)

	return rc, remResp
}

func (v *Verifier) VerifySig(m *big.Int, resp []*big.Int) bool {
	if len(resp) != 3 {
		return false
	}
	remResp := resp[:len(resp)-1]
	chlg := resp[len(resp)-1]

	rv, remResp := v.RecoverCommitment(chlg, remResp)

	if len(remResp) != 0 {
		return false
	}

	// c = H(M, rv)
	c := Hash(v.Q, m, rv)
	if c.Cmp(chlg) != 0 {
		return false
	}
	return true

}
