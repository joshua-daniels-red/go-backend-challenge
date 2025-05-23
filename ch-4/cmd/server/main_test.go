package main

import (
	"os"
	"testing"
)

func TestMainRunsWithoutError(t *testing.T) {
	// Write a temporary config.json
	config := `{
		"port": "7999",
		"stream_url": "https://example.com",
		"storage": "in-memory",
		"jwt_secret": "testsecret"
	}`
	err := os.WriteFile("config.json", []byte(config), 0644)
	if err != nil {
		t.Fatalf("failed to write config.json: %v", err)
	}
	defer os.Remove("config.json")

	// Run main in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main panicked: %v", r)
			}
		}()
		main()
	}()
}
