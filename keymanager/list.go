package keymanager

import (
	"encoding/json"
	"fmt"
    "os"
)

// ListKeys lists all SSH keys stored in the JSON file
func ListKeys(configPath string) error {
	// Load the config
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}
	
	// Mount the encrypted container
	if err := MountContainer(config); err != nil {
		return fmt.Errorf("failed to mount container: %v", err)
	}
	defer UnmountContainer()

	// Read the JSON file
	data, err := os.ReadFile(keysJson)
	if err != nil {
		return fmt.Errorf("failed to read keys JSON file: %v", err)
	}

	// Parse the JSON data
	var keys []KeyRecord
	if err := json.Unmarshal(data, &keys); err != nil {
		return fmt.Errorf("failed to parse keys JSON data: %v", err)
	}

	// Print the keys
	for _, key := range keys {
		fmt.Printf("Name: %s\nPublic Key: %s\n\n", key.Name, key.PubKey)
	}

	return nil
}
