package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type mockUserStore struct{}

func (m *mockUserStore) ValidateCredentials(username, password string) bool {
	return username == "admin" && password == "password123"
}

func TestLoginHandler_Success(t *testing.T) {
	store := &mockUserStore{}
	secret := "testsecret"
	handler := LoginHandler(store, secret)

	reqBody, _ := json.Marshal(LoginRequest{
		Username: "admin",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp LoginResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	parsedToken, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	store := &mockUserStore{}
	secret := "testsecret"
	handler := LoginHandler(store, secret)

	reqBody, _ := json.Marshal(LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginHandler_BadRequest(t *testing.T) {
	store := &mockUserStore{}
	secret := "testsecret"
	handler := LoginHandler(store, secret)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("not-json")))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
