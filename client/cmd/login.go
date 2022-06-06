package cmd

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"strconv"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login ",
	Run: func(cmd *cobra.Command, args []string) {
		server := cmd.Flag("server").Value.String()
		user := cmd.Flag("username").Value.String()
		password, err := strconv.Atoi(cmd.Flag("password").Value.String())
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		prover := zkp.NewProver(int64(password))
		request, err := prover.CreateAuthenticationCommits()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		// construct login request
		authRequest := &zkp_pb.AuthRequest{
			User: user,
			Commits: []*zkp_pb.AuthRequest_Commits{
				&zkp_pb.AuthRequest_Commits{
					R1: request.C1.Bytes(),
					R2: request.C2.Bytes(),
				},
			},
		}

		// connect to server
		conn, err := net.Dial("tcp", server)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer conn.Close()

		// send request
		if err := zkp.SendMessage(conn, authRequest); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		// wait for response
		msg, err := zkp.ReadMessage(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		challengeResponse, ok := msg.(*zkp_pb.ChallengeResponse)
		if !ok {
			fmt.Println("Error: wrong auth response")
			return
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
			return
		}

		// wait for response
		msg, err = zkp.ReadMessage(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		authResponse, ok := msg.(*zkp_pb.AuthResponse)
		if !ok {
			fmt.Println("Error: wrong auth response")
			return
		}

		if authResponse.Result {
			fmt.Println("Login successful")
		} else {
			if authResponse.Error == "" {
				fmt.Println("Error: wrong user name or password")
			} else {
				fmt.Printf("Error: %s\n", authResponse.Error)
			}
		}
	},
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()

	loginCmd.PersistentFlags().StringP("username", "u", "", "username")
	viper.BindPFlag("username", flags.Lookup("username"))
	// todo: set required field and validate input

	loginCmd.PersistentFlags().Int16P("password", "p", 0, "password")
	viper.BindPFlag("password", flags.Lookup("password"))
	// todo: set required field and validate input

	rootCmd.AddCommand(loginCmd)
}
