package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "barn",
	Short: "A CLI tool to manage SSH keys stored in an encrypted VeraCrypt container",
	Long:  `A CLI tool to manage SSH keys stored in an encrypted VeraCrypt container.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "config file (default is config.json)")
}