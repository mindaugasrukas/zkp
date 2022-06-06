package app

import (
	"errors"
	"log"
	"math/big"
	"net"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

func (s *Server) serveAuth(conn net.Conn, authRequest *zkp_pb.AuthRequest) error {
	if len(authRequest.GetCommits()) == 0 {
		// todo: wrong request
		return WrongRequestError
	}

	c := authRequest.GetCommits()[0]

	var r1, r2 big.Int
	r1.SetBytes(c.GetR1())
	r2.SetBytes(c.GetR2())

	user := zkp.UUID(authRequest.GetUser())
	auth := zkp.Commits{
		C1: &r1,
		C2: &r2,
	}
	log.Printf("r1=%v, r2=%v", &r1, &r2)
	return s.authenticate(conn, user, &auth)
}

func (s *Server) authenticate(connection net.Conn, user zkp.UUID, authRequest *zkp.Commits) error {
	// Get the user data
	userCommits, err := s.registry.Get(user)
	if err != nil {
		return err
	}

	// Send the challenge
	challenge, err := s.Verifier.CreateAuthenticationChallenge()
	if err != nil {
		return err
	}
	challengeResponse := &zkp_pb.ChallengeResponse{
		Challenge: (*big.Int)(challenge).Bytes(),
	}
	if err := zkp.SendMessage(connection, challengeResponse); err != nil {
		return err
	}

	// Verify the answer
	msg, err := zkp.ReadMessage(connection)
	if err != nil {
		return err
	}
	answerRequest, ok := msg.(*zkp_pb.AnswerRequest)
	if !ok {
		return errors.New("wrong auth answer")
	}

	var answer big.Int
	answer.SetBytes(answerRequest.GetAnswer())
	log.Print("answer = ", &answer)

	result := s.Verifier.VerifyAuthentication(userCommits, authRequest, &answer)

	// Send authentication results
	authResponse := &zkp_pb.AuthResponse{
		Result: result,
	}
	if err := zkp.SendMessage(connection, authResponse); err != nil {
		return err
	}

	return nil
}

