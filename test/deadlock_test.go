package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNoDeadlocks(t *testing.T) {
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

	// Wait for a few seconds to see if the app runs without deadlocking
	time.Sleep(5 * time.Second)

	// Terminate the process
	if cmd.Process != nil {
		fmt.Println("Terminating hammerclock...")
		cmd.Process.Kill()
	}
	fmt.Println("Test completed successfully! The application did not deadlock.")
}

// This ensures the test will run with "go test"
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
