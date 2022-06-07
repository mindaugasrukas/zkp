package app

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/mindaugasrukas/zkp_example/server/model"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
)

func (s *Server) serveAuth(conn net.Conn, authRequest *zkp_pb.AuthRequest) error {
	if len(authRequest.GetCommits()) == 0 {
		// todo: wrong request
		return WrongRequestError
	}

	user, auth := model.GetAuthentication(authRequest)
	if err := s.authenticate(conn, user, auth); err != nil {
		// send error response
		authResponse := &zkp_pb.AuthResponse{
			Result: false,
			Error:  err.Error(),
		}
		if err := zkp.SendMessage(conn, authResponse); err != nil {
			// log the error and continue
			fmt.Println(err.Error())
		}
		return err
	}

	return nil
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
		Challenge: (challenge).Bytes(),
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

	answer := model.GetAnswer(answerRequest)
	log.Print("answer = ", &answer)

	result := s.Verifier.VerifyAuthentication(userCommits, authRequest, challenge, answer)

	// Send authentication results
	authResponse := &zkp_pb.AuthResponse{
		Result: result,
	}
	if err := zkp.SendMessage(connection, authResponse); err != nil {
		return err
	}

	return nil
}
