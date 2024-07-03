package keymanager

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"sync"
    "syscall"
)

func lockMemory(data []byte) error {
	return unix.Mlock(data)
}

func unlockMemory(data []byte) error {
	return unix.Munlock(data)
}

func zeroMemory(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// Setup initializes the encrypted container and key file
func Setup(configPath, containerPath, keyfilePath, gpgKeyID string) error {
    keyfileSize := 32 * 1024 * 1024 // 32 MB

    // Generate a large random key file in memory
    randomData := make([]byte, keyfileSize)
    if err := lockMemory(randomData); err != nil {
        return err
    }
    defer unlockMemory(randomData)
    defer zeroMemory(randomData)

    if _, err := rand.Read(randomData); err != nil {
        return err
    }

    // Extract key and password from the random data
    actualKey, password, err := extractKeyAndPassword(randomData)
    if err != nil {
        return err
    }

    var wg sync.WaitGroup
    wg.Add(2)

    // Encrypt the key file with GPG in parallel
    go func() {
        defer wg.Done()
        // Create a pipe for passing the random data to GPG
        randomR, randomW, err := os.Pipe()
        if err != nil {
            panic(err)
        }

        // Write the random data to the pipe in a separate goroutine
        go func() {
            defer randomW.Close()
            if _, err := randomW.Write(randomData); err != nil {
                panic(err)
            }
        }()

        // Encrypt the data with GPG
        encryptCmd := exec.Command("gpg", "--yes", "--output", keyfilePath+".gpg", "--encrypt", "--recipient", gpgKeyID)
        encryptCmd.Stdin = randomR
        if err := encryptCmd.Run(); err != nil {
            panic(err)
        }

        // Clean up the pipe
        if err := randomR.Close(); err != nil {
            panic(err)
        }
    }()

    // Create a VeraCrypt container in parallel
    go func() {
        defer wg.Done()

        // Create a named pipe for passing the key to VeraCrypt
    	keyfilePipe := "mykeyfilepipe"
    	// Remove the named pipe if it already exists
    	if _, err := os.Stat(keyfilePipe); err == nil {
    		os.Remove(keyfilePipe)
    	}
    	err = syscall.Mkfifo(keyfilePipe, 0600)
    	if err != nil {
			panic(err)
    	}
    	defer os.Remove(keyfilePipe)
    
//    	// Write the actual key to the named pipe in a separate goroutine
//    	go func() {
//    		keyfile, err := os.OpenFile(keyfilePipe, os.O_WRONLY, 0600)
//    		if err != nil {
//    			panic(fmt.Errorf("error opening named pipe for writing: %v", err))
//    		}
//    		defer keyfile.Close()
//    
//    		if _, err := keyfile.Write(actualKey); err != nil {
//    			panic(err)
//    		}
//    	}()
		
		var wg sync.WaitGroup
		wg.Add(1)
		go WriteKeyToPipe(actualKey, keyfilePipe, &wg)
	
		// Wait for the goroutine to finish
		wg.Wait()

        // Prepare the VeraCrypt command
        cmd := exec.Command(veracryptPath, "--text", "--non-interactive", "--create", containerPath, "--size=32M", "--encryption=AES", "--hash=SHA-512", "--filesystem=fat", "--keyfiles="+keyfilePipe, "--password="+string(password))

        // Start the VeraCrypt command
        if err := cmd.Start(); err != nil {
            fmt.Printf("Command error: %v\n", err)
            panic(err)
        }

        // Wait for the VeraCrypt command to finish
        if err := cmd.Wait(); err != nil {
            fmt.Printf("Command error: %v\n", err)
            panic(err)
        }
    }()

    wg.Wait()

    // Save the config
    config := &Config{
        ContainerPath: containerPath,
        KeyfilePath:   keyfilePath + ".gpg",
        GPGKeyID:      gpgKeyID,
        }
        if err := SaveConfig(config, configPath); err != nil {
            return err
        }

        fmt.Println("Setup completed successfully.")
        return nil
}
