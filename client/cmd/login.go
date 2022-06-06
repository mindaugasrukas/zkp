package cmd

import (
	"fmt"
	"strconv"

	"github.com/mindaugasrukas/zkp_example/client/app"
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

		client := app.NewClient(server)
		if err = client.Login(user, password); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	},
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()

	loginCmd.PersistentFlags().StringP("username", "u", viper.GetString("USER"), "username (env: USER)")
	viper.BindPFlag("username", flags.Lookup("username"))
	// todo: set required field and validate input

	loginCmd.PersistentFlags().Int16P("password", "p", int16(viper.GetInt("PASSWORD")), "password (env: PASSWORD)")
	viper.BindPFlag("password", flags.Lookup("password"))
	// todo: set required field and validate input

	rootCmd.AddCommand(loginCmd)
}
