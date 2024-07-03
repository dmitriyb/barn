package keymanager

import (
	"encoding/json"
	"fmt"
	"os"
)

// KeyRecord represents an SSH key record in the JSON file
type KeyRecord struct {
	Name     string `json:"name"`
	PubKey   string `json:"pub_key"`
	FilePath string `json:"file_path"`
}

// ReadKeysFromJSON reads and parses the keys from the JSON file
func ReadKeysFromJSON(keysJson string) ([]KeyRecord, error) {
	// Read the keys JSON file
	data, err := os.ReadFile(keysJson)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys JSON file: %v", err)
	}

	// Parse the JSON data
	var keys []KeyRecord
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse keys JSON data: %v", err)
	}

	return keys, nil
}

// FindKeyByName finds a key by name from the JSON file
func FindKeyByName(keysJson, keyName string) (*KeyRecord, error) {
	// Read the keys from the JSON file
	keys, err := ReadKeysFromJSON(keysJson)
	if err != nil {
		return nil, fmt.Errorf("unable to read keys from keys JSON: %v", err)
	}

	// Find the key by name
	for _, k := range keys {
		if k.Name == keyName {
			return &k, nil
		}
	}

	return nil, fmt.Errorf("key with name %s not found", keyName)
}

// AddRecordToKeys adds a new key record to the keys JSON file
func AddRecordToKeys(jsonPath string, record KeyRecord) error {
	// Read the JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Parse the keys JSON data
	var keys []KeyRecord
	if len(data) > 0 {
		if err := json.Unmarshal(data, &keys); err != nil {
			return fmt.Errorf("failed to parse JSON data: %v", err)
		}
	}

	// Add the new key to the list
	keys = append(keys, record)

	// Write the updated keys JSON data back to the file
	updatedData, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}
	if err := os.WriteFile(jsonPath, updatedData, 0600); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}

// DeleteRecordFromKeys deletes a key record from the keys JSON file by name
func DeleteRecordFromKeys(jsonPath, name string) error {
	// Read the JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Parse the JSON data
	var keys []KeyRecord
	if err := json.Unmarshal(data, &keys); err != nil {
		return fmt.Errorf("failed to parse JSON data: %v", err)
	}

	// Find and delete the key by name
	for i, key := range keys {
		if key.Name == name {
			keys = append(keys[:i], keys[i+1:]...)
			break
		}
	}

	// Write the updated keys JSON data back to the file
	updatedData, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}
	if err := os.WriteFile(jsonPath, updatedData, 0600); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}

// CheckAndCreateKeysStructure checks if the JSON file exists and creates it if it does not
func CheckAndCreateKeysStructure(jsonPath string) error {
	_, err := os.Stat(jsonPath)
	if os.IsNotExist(err) {
		// Create an empty keys JSON array
		emptyJSON := []KeyRecord{}

		data, err := json.MarshalIndent(emptyJSON, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal empty keys JSON data: %v", err)
		}

		if err := os.WriteFile(jsonPath, data, 0600); err != nil {
			return fmt.Errorf("failed to write keys JSON file: %v", err)
		}
		
		// Define the directory for the SSH keys inside the mounted container
		if err := os.MkdirAll(keysPath, 0700); err != nil {
			return fmt.Errorf("failed to create key directory: %v", err)
		}
		
	} else if err != nil {
		return fmt.Errorf("failed to check keys JSON file: %v", err)
	}

	return nil
}
