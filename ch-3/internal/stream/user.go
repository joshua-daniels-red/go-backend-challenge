package stream

import (
	"github.com/gocql/gocql"
)

type UserStore struct {
	session *gocql.Session
}

type CredentialValidator interface {
	ValidateCredentials(username, password string) bool
}

func NewUserStore(session *gocql.Session) *UserStore {
	return &UserStore{session: session}
}

func (us *UserStore) ValidateCredentials(username, password string) bool {
	var storedPassword string
	if err := us.session.Query(`SELECT password FROM goanalytics.users WHERE username = ?`, username).
		Scan(&storedPassword); err != nil {
		return false
	}
	return storedPassword == password
}