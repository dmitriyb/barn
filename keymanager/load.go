package keymanager

import (
	"fmt"
	"os/exec"
)

// LoadKey loads an SSH key into the SSH agent by its name
func LoadKey(configPath, keyName string) error {
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
	
	keyToLoad, err := FindKeyByName(keysJson, keyName)
    if err != nil {
        return err
    }

	// Load the key into the SSH agent
	cmd := exec.Command("ssh-add", keyToLoad.FilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %v, output: %s", err, output)
	}

	fmt.Printf("Key with name %s and public key %s was added to the agent\n", keyToLoad.Name, keyToLoad.PubKey)
	return nil
}
