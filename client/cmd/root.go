package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().StringP("server", "s", "", "server URL")
	_ = rootCmd.MarkPersistentFlagRequired("server")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
