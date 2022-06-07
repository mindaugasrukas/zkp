package app

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/mindaugasrukas/zkp_example/client/model"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

var (
	UnknownResponseError = errors.New("unknown response")
	WrongResponseError   = errors.New("wrong response")
)

type (
	// Prover interface
	Prover interface {
		CreateAuthenticationCommits() (*zkp.Commits, error)
		ProveAuthentication(challenge *big.Int) (answer *big.Int)
	}

	Client struct {
		serverAddr string
		// Pluggable ZKP prover
		prover Prover
	}
)

// NewClient returns a new client instance
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
	}
}

// Register user to the server
func (c *Client) Register(user string, password int) error {
	// connect to server
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	defer conn.Close()

	prover := zkp.NewProver(int64(password))
	commits, err := prover.CreateRegisterCommits()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	log.Printf("y1=%v, y2=%v", commits.C1, commits.C2)

	// construct registration request
	request := &zkp_pb.RegisterRequest{
		User: user,
		Commits: []*zkp_pb.RegisterRequest_Commits{
			{
				Y1: commits.C1.Bytes(),
				Y2: commits.C2.Bytes(),
			},
		},
	}

	// send request
	if err := zkp.SendMessage(conn, request); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return c.ProcessResponse(conn)
}

// ProcessRegistrationResults ...
func (c *Client) ProcessRegistrationResults(registerResponse *zkp_pb.RegisterResponse) error {
	if registerResponse.Result {
		fmt.Println("Registration successful")
	} else {
		fmt.Printf("Error: %s\n", registerResponse.Error)
	}
	return nil
}

// Login user against the server
func (c *Client) Login(user string, password int) error {
	// connect to server
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	defer conn.Close()

	c.prover = zkp.NewProver(int64(password))
	request, err := c.prover.CreateAuthenticationCommits()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	// construct login request
	authRequest := &zkp_pb.AuthRequest{
		User: user,
		Commits: []*zkp_pb.AuthRequest_Commits{
			{
				R1: request.C1.Bytes(),
				R2: request.C2.Bytes(),
			},
		},
	}

	// send request
	if err := zkp.SendMessage(conn, authRequest); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return c.ProcessResponse(conn)
}

// ProcessChallenge returns answer to the server
func (c *Client) ProcessChallenge(conn net.Conn, challengeResponse *zkp_pb.ChallengeResponse) error {
	// construct answer request
	challenge := model.GetChallenge(challengeResponse)
	log.Print("challenge = ", challenge)
	answer := c.prover.ProveAuthentication(challenge)
	log.Print("answer = ", answer)
	answerRequest := &zkp_pb.AnswerRequest{
		Answer: (answer).Bytes(),
	}

	// send request
	if err := zkp.SendMessage(conn, answerRequest); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return c.ProcessResponse(conn)
}

// ProcessAuthResults ...
func (c *Client) ProcessAuthResults(authResponse *zkp_pb.AuthResponse) error {
	if authResponse.Result {
		fmt.Println("Login successful")
		for {
			// run infinite loop
		}
	} else {
		if authResponse.Error == "" {
			fmt.Println("Error: wrong user name or password")
		} else {
			fmt.Printf("Error: %s\n", authResponse.Error)
		}
	}

	return nil
}

// ProcessResponse will wait and process server response
// Very naive command processor
func (c *Client) ProcessResponse(conn net.Conn) error {
	msg, err := zkp.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	switch string(msg.ProtoReflect().Descriptor().Name()) {
	case "RegisterResponse":
		registerResponse, ok := msg.(*zkp_pb.RegisterResponse)
		if !ok {
			return WrongResponseError
		}
		return c.ProcessRegistrationResults(registerResponse)
	case "ChallengeResponse":
		challengeResponse, ok := msg.(*zkp_pb.ChallengeResponse)
		if !ok {
			return WrongResponseError
		}
		return c.ProcessChallenge(conn, challengeResponse)
	case "AuthResponse":
		authResponse, ok := msg.(*zkp_pb.AuthResponse)
		if !ok {
			return WrongResponseError
		}
		return c.ProcessAuthResults(authResponse)
	}

	return UnknownResponseError
}
