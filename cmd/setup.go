package cmd

import (
	"fmt"
	"os"
	
	"github.com/dmitriyb/barn/keymanager"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the encrypted container and key file",
	Long:  `Setup the encrypted container and key file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		containerPath, _ := cmd.Flags().GetString("container")
		keyfilePath, _ := cmd.Flags().GetString("keyfile")
		gpgKeyID, _ := cmd.Flags().GetString("gpg-key")

		err := keymanager.Setup(configPath, containerPath, keyfilePath, gpgKeyID)
		if err != nil {
			fmt.Printf("Error setting up: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringP("container", "C", "", "Path to VeraCrypt container")
	setupCmd.Flags().StringP("keyfile", "K", "", "Path to keyfile")
	setupCmd.Flags().StringP("gpg-key", "G", "", "GPG key ID")
	setupCmd.MarkFlagRequired("container")
	setupCmd.MarkFlagRequired("keyfile")
	setupCmd.MarkFlagRequired("gpg-key")
}