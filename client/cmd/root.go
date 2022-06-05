package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Dummy client application",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
)

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()
	flags.StringP("server", "s", viper.GetString("SERVER"), "server URL (env: SERVER)")
	viper.BindPFlag("server", flags.Lookup("server"))
	// todo: set required field and validate input
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
