package stream

import (
	"github.com/gocql/gocql"
)

type UserStore interface {
	ValidateCredentials(username, password string) bool
}


type CassandraUserStore struct {
	session *gocql.Session
}

func NewUserStore(session *gocql.Session) *CassandraUserStore {
	return &CassandraUserStore{session: session}
}

func (us *CassandraUserStore) ValidateCredentials(username, password string) bool {
	var storedPassword string
	if err := us.session.Query(`SELECT password FROM goanalytics.users WHERE username = ?`, username).
		Scan(&storedPassword); err != nil {
		return false
	}
	return storedPassword == password
}

type InMemoryUserStore struct{}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{}
}

func (s *InMemoryUserStore) ValidateCredentials(username, password string) bool {
	return username == "admin" && password == "admin"
}
