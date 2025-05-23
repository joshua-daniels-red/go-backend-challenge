package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockUserStore is a mock implementation for testing.
type mockUserStore struct{}

func (m *mockUserStore) ValidateCredentials(username, password string) bool {
	return username == "admin" && password == "password123"
}

func TestValidateCredentials(t *testing.T) {
	store := &mockUserStore{}

	assert.True(t, store.ValidateCredentials("admin", "password123"))
	assert.False(t, store.ValidateCredentials("admin", "wrongpass"))
	assert.False(t, store.ValidateCredentials("user", "password123"))
}
