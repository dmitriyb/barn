package keymanager

import (
    "fmt"
	"os"
)

// RemoveKey removes a specified SSH key from the encrypted storage
func RemoveKey(configPath, keyName string) error {
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
	
	keyToDelete, err := FindKeyByName(keysJson, keyName)
    if err != nil {
        return err
    }

	// Remove the key files from the file system
	if err := os.Remove(keyToDelete.FilePath); err != nil {
		return fmt.Errorf("failed to remove key file: %v", err)
	}
	if err := os.Remove(keyToDelete.FilePath + ".pub"); err != nil {
		return fmt.Errorf("failed to remove public key file: %v", err)
	}

	// Remove the key record from the JSON file
	if err := DeleteRecordFromKeys(keysJson, keyName); err != nil {
		return fmt.Errorf("failed to delete record from JSON file: %v", err)
	}

	fmt.Printf("Key with name %s was removed from the file system and JSON file\n", keyToDelete.Name)
	return nil
}
