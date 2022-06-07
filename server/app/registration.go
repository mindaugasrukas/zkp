package app

import (
	"fmt"
	"net"

	"github.com/mindaugasrukas/zkp_example/server/model"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

func (s *Server) serveRegistration(conn net.Conn, registerRequest *zkp_pb.RegisterRequest) error {
	if len(registerRequest.GetCommits()) == 0 {
		// todo: wrong request
		return WrongRequestError
	}

	user, commits := model.GetRegistration(registerRequest)
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
