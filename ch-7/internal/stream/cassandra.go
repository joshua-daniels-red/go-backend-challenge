
package stream

import (
	"log"
)

type Session interface {
	Query(stmt string, values ...interface{}) Query
}

type Query interface {
	Exec() error
	Iter() Iter
}

type Iter interface {
	Scan(dest ...interface{}) bool
	Close() error
}

type CassandraStats struct {
	session Session
}

func NewCassandraStats(session Session) *CassandraStats {
	return &CassandraStats{session: session}
}

func (c *CassandraStats) Record(event Event) {
	if err := c.session.Query(`
		UPDATE stats_by_domain SET count = count + 1 WHERE domain = ?
	`, event.Domain).Exec(); err != nil {
		log.Printf("failed to update stats_by_domain: %v", err)
	}

	if err := c.session.Query(`
		UPDATE stats_by_user SET count = count + 1 WHERE user = ?
	`, event.User).Exec(); err != nil {
		log.Printf("failed to update stats_by_user: %v", err)
	}
}

func (c *CassandraStats) RecordMany(events []Event) {
	for _, event := range events {
		if err := c.session.Query(`
			UPDATE stats_by_domain SET count = count + 1 WHERE domain = ?
		`, event.Domain).Exec(); err != nil {
			log.Printf("failed to update stats_by_domain: %v", err)
		}

		if err := c.session.Query(`
			UPDATE stats_by_user SET count = count + 1 WHERE user = ?
		`, event.User).Exec(); err != nil {
			log.Printf("failed to update stats_by_user: %v", err)
		}
	}
}


func (c *CassandraStats) GetSnapshot() StatsSnapshot {
	snapshot := StatsSnapshot{
		ByDomain: make(map[string]int),
		ByUser:   make(map[string]int),
	}

	iter := c.session.Query(`SELECT domain, count FROM stats_by_domain`).Iter()
	var domain string
	var count int
	for iter.Scan(&domain, &count) {
		snapshot.ByDomain[domain] = count
	}
	if err := iter.Close(); err != nil {
		log.Printf("error closing stats_by_domain iterator: %v", err)
	}

	iter = c.session.Query(`SELECT user, count FROM stats_by_user`).Iter()
	var user string
	for iter.Scan(&user, &count) {
		snapshot.ByUser[user] = count
	}
	if err := iter.Close(); err != nil {
		log.Printf("error closing stats_by_user iterator: %v", err)
	}

	return snapshot
}
