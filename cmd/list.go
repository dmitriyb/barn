package cmd

import (
	"fmt"
	"os"
	
    "github.com/dmitriyb/barn/keymanager"
    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved SSH keys stored inside the encrypted storage",
	Long:  `List all saved SSH keys stored inside the encrypted storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")

		err := keymanager.ListKeys(configPath)
		if err != nil {
			fmt.Printf("Error listing keys: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
