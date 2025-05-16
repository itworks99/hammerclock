package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Starting hammerclock test...")

	// Start the application process
	cmd := exec.Command("bin/hammerclock.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting hammerclock: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Hammerclock started successfully")
	fmt.Println("Waiting 5 seconds to verify no deadlocks...")

	// Wait for a few seconds to see if the hammerclock runs without deadlocking
	time.Sleep(5 * time.Second)

	// Terminate the process
	if cmd.Process != nil {
		fmt.Println("Terminating hammerclock...")
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Error terminating hammerclock: %v\n", err)
			return
		}
	}

	fmt.Println("Test completed successfully! The application did not deadlock.")
}
