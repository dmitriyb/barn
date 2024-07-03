package cmd

import (
    "fmt"
	"os"
	
    "github.com/dmitriyb/barn/keymanager"
    "github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load an SSH key from the encrypted storage into the SSH agent",
	Long:  `Load an SSH key from the encrypted storage into the SSH agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		keyName, _ := cmd.Flags().GetString("keyName")
		
		err := keymanager.LoadKey(configPath, keyName)
		if err != nil {
			fmt.Printf("Error loading key: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringP("keyName", "k", "", "Name of the key")
	loadCmd.MarkFlagRequired("keyName")
}
