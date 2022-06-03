package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register ",
	Run: func(cmd *cobra.Command, args []string) {
		server := cmd.Flag("server").Value.String()
		user := cmd.Flag("username").Value.String()
		password, err := strconv.Atoi(cmd.Flag("password").Value.String())
		if err != nil {
			fmt.Println(err)
			return
		}
		prover := zkp.NewProver(int64(password))
		commits := prover.CreateRegisterCommits()

		// construct request
		request := &zkp_pb.RegisterRequest{
			User: user,
			Commits: []*zkp_pb.RegisterRequest_Commits{
				&zkp_pb.RegisterRequest_Commits{
					Y1: commits.Y1.Bytes(),
					Y2: commits.Y2.Bytes(),
				},
			},
		}

		// connect to server
		connection, err := net.Dial("tcp", server)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer connection.Close()

		// send request
		out, err := proto.Marshal(request)
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err = connection.Write(out); err != nil {
			fmt.Println(err)
			return
		}

		// wait for response
		// todo:
	},
}

func init() {
	registerCmd.PersistentFlags().StringP("username", "u", "", "username")
	_ = registerCmd.MarkPersistentFlagRequired("username")
	registerCmd.PersistentFlags().Int16P("password", "p", 0, "password")
	_ = registerCmd.MarkPersistentFlagRequired("password")

	rootCmd.AddCommand(registerCmd)
}