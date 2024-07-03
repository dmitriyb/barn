package keymanager

import (
	"bufio"
	"bytes"
    "encoding/json"
    "fmt"
	"io"
	"os"
	"os/exec"
    "strings"
    "sync"
    "syscall"
    "time"
)

const veracryptPath = "/usr/local/bin/veracrypt" // Update this to the full path of VeraCrypt on your system
const mountPoint = "/Volumes/barn"
const keysJson = mountPoint + "/keys.json"

const keysPath = mountPoint + "/keys"

type Config struct {
	ContainerPath string `json:"container_path"`
	KeyfilePath   string `json:"keyfile_path"`
	GPGKeyID      string `json:"gpg_key_id"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to a JSON file
func SaveConfig(config *Config, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func WriteKeyToPipe(actualKey []byte, keyfilePipe string, wg *sync.WaitGroup) {
	defer wg.Done()

	keyfile, err := os.OpenFile(keyfilePipe, os.O_WRONLY, 0600)
	if err != nil {
		panic(fmt.Errorf("error opening named pipe for writing: %v", err))
	}
	defer keyfile.Close()

	if _, err := keyfile.Write(actualKey); err != nil {
		panic(err)
	}
}

// MountContainer mounts the VeraCrypt container using the decrypted key file and password
func MountContainer(config *Config) error {
	// Ensure GPG_TTY environment variable is set for interactive pinentry
	tty := os.Getenv("GPG_TTY")
	if tty == "" {
		ttyBytes, err := exec.Command("tty").Output()
		if err != nil {
			return fmt.Errorf("failed to get tty: %v", err)
		}
		tty = strings.TrimSpace(string(ttyBytes))
		os.Setenv("GPG_TTY", tty)
	}

	// Decrypt the key file using GPG
	fmt.Println("Starting GPG decryption...")
	decryptCmd := exec.Command("gpg", "--decrypt", config.KeyfilePath)
	stdout, err := decryptCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := decryptCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}
	if err := decryptCmd.Start(); err != nil {
		return fmt.Errorf("failed to start GPG command: %v", err)
	}

	// Capture GPG stderr to handle potential PIN prompts
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("GPG stderr:", scanner.Text()) // Print the GPG stderr output
		}
	}()

	// Read the decrypted key file
	randomData, err := io.ReadAll(bufio.NewReader(stdout))
	if err != nil {
		return fmt.Errorf("failed to read decrypted data: %v", err)
	}

	// Extract key and password from the random data
	actualKey, password, err := extractKeyAndPassword(randomData)
	if err != nil {
		return fmt.Errorf("insufficient data for offsets: %v", err)
	}

	// Ensure the GPG command has completed
	if err := decryptCmd.Wait(); err != nil {
		return fmt.Errorf("GPG command failed: %v", err)
	}

	// Create a named pipe for passing the key to VeraCrypt
	keyfilePipe := "mykeyfilepipe"
	// Remove the named pipe if it already exists
	if _, err := os.Stat(keyfilePipe); err == nil {
		os.Remove(keyfilePipe)
	}
	err = syscall.Mkfifo(keyfilePipe, 0600)
	if err != nil {
		return fmt.Errorf("error creating named pipe: %v", err)
	}
	defer os.Remove(keyfilePipe)

//	// Write the actual key to the named pipe in a separate goroutine
//	go func() {
//		keyfile, err := os.OpenFile(keyfilePipe, os.O_WRONLY, 0600)
//		if err != nil {
//			panic(fmt.Errorf("error opening named pipe for writing: %v", err))
//		}
//		defer keyfile.Close()
//
//		if _, err := keyfile.Write(actualKey); err != nil {
//			panic(err)
//		}
//	}()
	
	var wg sync.WaitGroup
	wg.Add(1)
	go WriteKeyToPipe(actualKey, keyfilePipe, &wg)

	// Wait for the goroutine to finish
	wg.Wait()

	// Mount the VeraCrypt volume using the actual key from the named pipe
	mountCmd := exec.Command("veracrypt", "--text", "--non-interactive", "--password="+string(password), "--keyfiles="+keyfilePipe, "--mount", config.ContainerPath, "/Volumes/barn")
	var out bytes.Buffer
	mountCmd.Stdout = &out
	mountCmd.Stderr = &out
	if err := mountCmd.Run(); err != nil {
		fmt.Printf("Command error: %v\n", err)
		fmt.Printf("Output: %s\n", out.String())
		return err
	}

	fmt.Println("VeraCrypt volume mounted successfully at " + mountPoint)

    err = CheckAndCreateKeysStructure(keysJson)
    if err != nil {
        return err
    }
	
	return nil
}

// UnmountContainer unmounts the VeraCrypt container
func UnmountContainer() error {
	cmd := exec.Command("veracrypt", "--text", "--dismount", mountPoint)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("VeraCrypt volume unmounted successfully.")
	return nil
}
