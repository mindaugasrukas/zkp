package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"

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

	// Verifier interface
	Verifier interface {
		CreateAuthenticationChallenge() (challenge *big.Int, err error)
		VerifyAuthentication(commits *zkp.Commits, authRequest *zkp.Commits, answer *big.Int) bool
	}

	// Server application
	Server struct {
		// Pluggable storage
		registry Registry
		// Pluggable ZKP verifier
		Verifier Verifier
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
		Verifier: zkp.NewVerifier(),
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

	log.Print("Listening on ", l.Addr())

	// run infinite loop
	for {
		// todo: add a rate limiter
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
			fmt.Println()
		}(conn)
	}
}

func (s *Server) serve(conn net.Conn) error {
	msg, err := zkp.ReadMessage(conn)
	if err != nil {
		return err
	}

	switch string(msg.ProtoReflect().Descriptor().Name()) {
	case "RegisterRequest":
		registerRequest, ok := msg.(*zkp_pb.RegisterRequest)
		if !ok {
			return WrongRequestError
		}
		return s.serveRegistration(conn, registerRequest)
	case "AuthRequest":
		authRequest, ok := msg.(*zkp_pb.AuthRequest)
		if !ok {
			return WrongRequestError
		}
		return s.serveAuth(conn, authRequest)
	}

	return UnknownRequestError
}
