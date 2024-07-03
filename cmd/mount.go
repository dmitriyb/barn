package cmd

import (
	"fmt"
	"os"
	
	"github.com/dmitriyb/barn/keymanager"
	"github.com/spf13/cobra"
)

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount a VeraCrypt volume",
	Long:  `Mount the VeraCrypt volume using the encrypted key file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")

		config, err := keymanager.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		err = keymanager.MountContainer(config)
		if err != nil {
			fmt.Printf("Error mounting VeraCrypt volume: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)
}
