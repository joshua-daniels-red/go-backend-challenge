package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func setupTempConfig(t *testing.T, port string) func() {
	t.Helper()
	configContent := fmt.Sprintf(`{
        "port": "%s",
        "stream_url": "https://example.com",
        "storage": "in-memory",
        "jwt_secret": "testsecret"
    }`, port)
	err := os.WriteFile("config.json", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write temp config.json: %v", err)
	}
	return func() {
		os.Remove("config.json")
	}
}

func TestMain_SuccessfulStartupAndGracefulShutdown(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	originalFlags := log.Flags()
	log.SetFlags(0) 
	defer func() {
		log.SetOutput(os.Stderr) 
		log.SetFlags(originalFlags)
	}()

	
	testPort := "8999" 
	cleanup := setupTempConfig(t, testPort)
	defer cleanup()

	mainDone := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main panicked: %v", r)
			}
			close(mainDone)
		}()
		main()
	}()

	time.Sleep(200 * time.Millisecond)

	client := http.Client{Timeout: 1 * time.Second}
	_, err := client.Get("http://localhost:" + testPort + "/some-test-endpoint") 
	if err == nil {
		log.Println("Test: Server responded to GET request before shutdown.")
	} else {
		log.Printf("Test: Server did not respond on port %s (err: %v). This might be okay if no routes are expected.", testPort, err)
	}


	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find current process: %v", err)
	}
	if err := p.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("Failed to send SIGINT: %v", err)
	}

	select {
	case <-mainDone:
	case <-time.After(7 * time.Second): 
		t.Fatal("main function did not terminate in time after SIGINT")
	}

	logs := logBuffer.String()
	if !strings.Contains(logs, "HTTP server listening on :"+testPort) {
		t.Errorf("Expected log message 'HTTP server listening on :%s', got logs:\n%s", testPort, logs)
	}
	if !strings.Contains(logs, "Shutting down server...") {
		t.Errorf("Expected log message 'Shutting down server...', got logs:\n%s", logs)
	}
	if !strings.Contains(logs, "Server exited properly") {
		t.Errorf("Expected log message 'Server exited properly', got logs:\n%s", logs)
	}

	_, err = client.Get("http://localhost:" + testPort)
	if err == nil {
		t.Errorf("Server still responding after expected shutdown on port %s", testPort)
	} else {
		log.Println("Test: Server is not responding after shutdown, as expected.")
	}
}
