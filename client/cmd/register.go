package cmd

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register ",
	Run: func(cmd *cobra.Command, args []string) {
		server := cmd.Flag("server").Value.String()
		user := cmd.Flag("username").Value.String()
		password, err := strconv.Atoi(cmd.Flag("password").Value.String())
		log.Print("password = ", password)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		prover := zkp.NewProver(int64(password))
		commits, err := prover.CreateRegisterCommits()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
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

		// connect to server
		conn, err := net.Dial("tcp", server)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer conn.Close()

		// send request
		if err := zkp.SendMessage(conn, request); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		// wait for response
		msg, err := zkp.ReadMessage(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		registerResponse, ok := msg.(*zkp_pb.RegisterResponse)
		if !ok {
			fmt.Println("Error: wrong registration response")
			return
		}

		if registerResponse.Result {
			fmt.Println("Registration successful")
		} else {
			fmt.Printf("Error: %s\n", registerResponse.Error)
		}
	},
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()

	registerCmd.PersistentFlags().StringP("username", "u", "", "username")
	viper.BindPFlag("username", flags.Lookup("username"))
	// todo: set required field and validate input

	registerCmd.PersistentFlags().Int16P("password", "p", 0, "password")
	viper.BindPFlag("password", flags.Lookup("password"))
	// todo: set required field and validate input

	rootCmd.AddCommand(registerCmd)
}
