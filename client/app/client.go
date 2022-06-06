package app

import (
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

type Client struct {
	serverAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
	}
}

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

	// wait for response
	msg, err := zkp.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	registerResponse, ok := msg.(*zkp_pb.RegisterResponse)
	if !ok {
		fmt.Println("Error: wrong registration response")
		return err
	}

	if registerResponse.Result {
		fmt.Println("Registration successful")
	} else {
		fmt.Printf("Error: %s\n", registerResponse.Error)
	}

	return nil
}

func (c *Client) Login(user string, password int) error {
	// connect to server
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	defer conn.Close()


	prover := zkp.NewProver(int64(password))
	request, err := prover.CreateAuthenticationCommits()
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

	// wait for response
	msg, err := zkp.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	challengeResponse, ok := msg.(*zkp_pb.ChallengeResponse)
	if !ok {
		fmt.Println("Error: wrong auth response")
		return err
	}

	// construct answer request
	var challenge big.Int
	challenge.SetBytes(challengeResponse.GetChallenge())
	log.Print("challenge = ", &challenge)
	answer := prover.ProveAuthentication(&challenge)
	log.Print("answer = ", answer)
	answerRequest := &zkp_pb.AnswerRequest{
		Answer: (*big.Int)(answer).Bytes(),
	}

	// send request
	if err := zkp.SendMessage(conn, answerRequest); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	// wait for response
	msg, err = zkp.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	authResponse, ok := msg.(*zkp_pb.AuthResponse)
	if !ok {
		fmt.Println("Error: wrong auth response")
		return err
	}

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
