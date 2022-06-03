package main

import (
	"errors"
	"fmt"
	"io"
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

var (
	WrongRequestError   = errors.New("wrong request")
	UnknownRequestError = errors.New("unknown request")
)

// NewServer returns a new server instance
func NewServer() *Server {
	return &Server{
		registry: store.NewInMemoryStore(),
		verifier: zkp.NewVerifier(),
	}
}

// Run starts the server
func (s *Server) Run(port string) {
	l, err := net.Listen("tcp", ":" + port)
	if err != nil {
		// Can't start - panic
		panic(err.Error())
	}
	defer l.Close()

	// run infinite loop
	for {
		conn, err := l.Accept()
		if err != nil {
			// log the error and continue
			fmt.Println(err.Error())
		}

		go func(conn net.Conn) {
			defer conn.Close()
			if err := s.serve(conn); err != nil {
				// log the error and continue
				fmt.Println(err.Error())
			}
		}(conn)
	}
}

func (s *Server) serve(conn net.Conn) error {
	in := make([]byte, 0, 10240)
	tmp := make([]byte, 4096)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		in = append(in, tmp[:n]...)
	}
	registerRequest := &zkp_pb.RegisterRequest{}
	if err := proto.Unmarshal(in, registerRequest); err == nil {
		return s.serveRegistration(registerRequest)
	}
	authRequest := &zkp_pb.AuthRequest{}
	if err := proto.Unmarshal(in, authRequest); err == nil {
		return s.serveAuth(conn, authRequest)
	}
	return UnknownRequestError
}

func (s *Server) serveRegistration(registerRequest *zkp_pb.RegisterRequest) error {
	if len(registerRequest.GetCommits()) == 0 {
		// todo: wrong request
		return WrongRequestError
	}

	c := registerRequest.GetCommits()[0]

	var y1, y2 big.Int
	y1.SetBytes(c.GetY1())
	y2.SetBytes(c.GetY2())

	commits := &zkp.Commits{
		Y1: &y1,
		Y2: &y2,
	}
	user := zkp.UUID(registerRequest.GetUser())
	if err := s.Register(user, commits); err != nil {
		// todo: return status to the client
		return fmt.Errorf("fail to register user %q: %v", user, err)
	}

	// todo: return status to the client
	fmt.Printf("registered new user %q\n", user)
	return nil
}

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
	auth := zkp.AuthenticationRequest{
		R1: &r1,
		R2: &r2,
	}
	return s.Authenticate(conn, user, auth)
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
