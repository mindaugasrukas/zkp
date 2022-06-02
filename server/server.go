package main

import (
	"math/big"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/mindaugasrukas/zkp_example/store"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

type (
	// Registry interface
	Registry interface {
		// Add user data to the registry
		Add(user zkp.UUID, commits *zkp.Commits) error
		// Get user data from the registry
		Get(user zkp.UUID) (*zkp.Commits, error)
	}

	Verifier interface {
		CreateAuthenticationChallenge(authRequest zkp.AuthenticationRequest) (zkp.Challenge, error)
		VerifyAuthentication(commits *zkp.Commits, authRequest zkp.AuthenticationRequest, answer zkp.Answer) bool
	}

	// Server application
	Server struct {
		registry Registry
		verifier Verifier
	}
)

// NewServer returns a new server instance
func NewServer() *Server {
	return &Server{
		registry: store.NewInMemoryStore(),
		verifier: zkp.NewVerifier(),
	}
}

// Register Registers a new user
func (s *Server) Register(user zkp.UUID, commits *zkp.Commits) error {
	return s.registry.Add(user, commits)
}

func (s *Server) Authenticate(connection net.Conn, user zkp.UUID, authRequest zkp.AuthenticationRequest) error {
	// Get the user data
	userCommits, err := s.registry.Get(user)
	if err != nil {
		return err
	}

	// Send the challenge
	challenge, err := s.verifier.CreateAuthenticationChallenge(authRequest)
	if err != nil {
		return err
	}
	challengeResponse := &zkp_pb.ChallengeResponse{
		Challenge: (*big.Int)(challenge).Bytes(),
	}
	out, err := proto.Marshal(challengeResponse)
	if err != nil {
		return err
	}
	if _, err = connection.Write(out); err != nil {
		return err
	}

	// Verify the answer
	in := make([]byte, 0)
	if _, err = connection.Read(in); err != nil {
		return err
	}
	answerRequest := &zkp_pb.AnswerRequest{}
	if err := proto.Unmarshal(in, answerRequest); err != nil {
		return err
	}

	var answer big.Int
	answer.SetBytes(answerRequest.GetAnswer())

	result := s.verifier.VerifyAuthentication(userCommits, authRequest, &answer)

	// Send authentication results
	authResponse := &zkp_pb.AuthResponse{
		Result: result,
	}
	out, err = proto.Marshal(authResponse)
	if err != nil {
		return err
	}
	if _, err = connection.Write(out); err != nil {
		return err
	}

	return nil
}
