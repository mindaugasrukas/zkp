package app

import (
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

func (s *Server) serveRegistration(conn net.Conn, registerRequest *zkp_pb.RegisterRequest) error {
	if len(registerRequest.GetCommits()) == 0 {
		// todo: wrong request
		return WrongRequestError
	}

	c := registerRequest.GetCommits()[0]

	var y1, y2 big.Int
	y1.SetBytes(c.GetY1())
	y2.SetBytes(c.GetY2())
	log.Printf("y1=%v, y2=%v", &y1, &y2)

	commits := &zkp.Commits{
		C1: &y1,
		C2: &y2,
	}
	user := zkp.UUID(registerRequest.GetUser())
	response := &zkp_pb.RegisterResponse{Result: true}

	if err := s.Register(user, commits); err != nil {
		response.Result = false
		response.Error = err.Error()
		if err := zkp.SendMessage(conn, response); err != nil {
			return err
		}
		return fmt.Errorf("fail to register user %q: %v", user, err)
	}

	fmt.Printf("registered new user %q\n", user)
	if err := zkp.SendMessage(conn, response); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Register Registers a new user
func (s *Server) Register(user zkp.UUID, commits *zkp.Commits) error {
	return s.registry.Add(user, commits)
}
