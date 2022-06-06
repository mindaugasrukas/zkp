package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mindaugasrukas/zkp_example/client/app"
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

		client := app.NewClient(server)
		if err = client.Register(user, password); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	},
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()

	registerCmd.PersistentFlags().StringP("username", "u", viper.GetString("USER"), "username")
	viper.BindPFlag("username", flags.Lookup("username"))
	// todo: set required field and validate input

	registerCmd.PersistentFlags().Int16P("password", "p", int16(viper.GetInt("PASSWORD")), "password")
	viper.BindPFlag("password", flags.Lookup("password"))
	// todo: set required field and validate input

	rootCmd.AddCommand(registerCmd)
}
