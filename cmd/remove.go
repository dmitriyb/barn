package cmd

import (
    "fmt"
	"os"
	
	"github.com/dmitriyb/barn/keymanager"
    "github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an SSH key from the encrypted storage",
	Long:  `Remove an SSH key from the encrypted storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		name, _ := cmd.Flags().GetString("keyName")
		
		// Call the RemoveKey function with the loaded config
		err := keymanager.RemoveKey(configPath, name)
		if err != nil {
			fmt.Printf("Error removing key: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringP("keyName", "k", "", "Name of the key")
	removeCmd.MarkFlagRequired("keyName")
}
