package cmd

import (
	"fmt"
	"os"
	
	"github.com/dmitriyb/barn/keymanager"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH key to the encrypted storage",
	Long:  `Add a new SSH key to the encrypted storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		name, _ := cmd.Flags().GetString("keyName")
		keyType, _ := cmd.Flags().GetString("type")
		comment, _ := cmd.Flags().GetString("comment")
		
		// Call the AddKey function with the loaded config
		err := keymanager.AddKey(configPath, name, keyType, comment)
		if err != nil {
			fmt.Printf("Error adding key: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("keyName", "k", "", "Name of the key")
	addCmd.Flags().StringP("type", "t", "ed25519", "Type of the key (e.g., rsa, ed25519)")
	addCmd.Flags().StringP("comment", "m", "", "Comment for the key")
	addCmd.MarkFlagRequired("keyName")
}
