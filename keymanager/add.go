package keymanager

import (
    "fmt"
    "os"
    "os/exec"
)

// AddKey adds a new SSH key to the encrypted storage
func AddKey(configPath, keyName, keyType, comment string) error {
	// Load the config
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}
	
	// Mount the encrypted container
	if err := MountContainer(config); err != nil {
		return fmt.Errorf("failed to mount container: %v", err)
	}
	defer func() {
		if err := UnmountContainer(); err != nil {
			fmt.Printf("failed to unmount container: %v", err)
		}
	}()
	
	// Define the path for the new SSH key
	keyPath := fmt.Sprintf("%s", keysPath + "/" + keyName)
	
	// Create the SSH key
	cmd := exec.Command("ssh-keygen", "-t", keyType, "-f", keyPath, "-C", comment, "-N", "")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %v, output: %s", err, output)
	}
	
	// Read the public key
	pubKeyPath := keyPath + ".pub"
	pubKeyData, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %v", err)
	}
	pubKey := string(pubKeyData)
	
	// Create the new key record
	record := KeyRecord{
		Name:     keyName,
		PubKey:   pubKey,
		FilePath: keyPath,
	}

	// Add the new key record to the JSON file
	if err := AddRecordToKeys(keysJson, record); err != nil {
		return fmt.Errorf("failed to add record to JSON file: %v", err)
	}

	fmt.Printf("Key with name %s and public key %s was created and added to the JSON file\n", keyName, pubKey)
	return nil
}
