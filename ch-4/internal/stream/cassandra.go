package stream

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

type CassandraStats struct {
	session *gocql.Session
}

func NewCassandraStats(host string) (*CassandraStats, error) {
	var session *gocql.Session
	var err error

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "goanalytics"
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 5

	// Retry loop
	for i := 0; i < 10; i++ {
		session, err = cluster.CreateSession()
		if err == nil {
			log.Printf("Connected to Cassandra on attempt %d", i+1)
			return &CassandraStats{session: session}, nil
		}

		log.Printf("Waiting for Cassandra (%d/10): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to Cassandra after retries: %w", err)
}

func (cs *CassandraStats) Record(ev ChangeEvent) {
	go func() {
		ctx := context.Background()

		if err := cs.session.Query(`UPDATE stats_summary SET total_messages = total_messages + 1 WHERE id = 'global'`).WithContext(ctx).Exec(); err != nil {
			log.Printf("error updating total_messages: %v", err)
		}

		if ev.Bot {
			if err := cs.session.Query(`UPDATE stats_summary SET bot_count = bot_count + 1 WHERE id = 'global'`).WithContext(ctx).Exec(); err != nil {
				log.Printf("error updating bot_count: %v", err)
			}
		} else {
			if err := cs.session.Query(`UPDATE stats_summary SET non_bot_count = non_bot_count + 1 WHERE id = 'global'`).WithContext(ctx).Exec(); err != nil {
				log.Printf("error updating non_bot_count: %v", err)
			}
		}

		if err := cs.session.Query(`INSERT INTO unique_users (username) VALUES (?)`, ev.User).WithContext(ctx).Exec(); err != nil {
			log.Printf("error inserting unique_user: %v", err)
		}

		if err := cs.session.Query(`UPDATE server_counts SET count = count + 1 WHERE server_url = ?`, ev.ServerURL).WithContext(ctx).Exec(); err != nil {
			log.Printf("error updating server_counts: %v", err)
		}
	}()
}

func (cs *CassandraStats) GetSnapshot() StatsSnapshot {
	ctx := context.Background()

	var total, bots, nonBots int
	if err := cs.session.Query(`SELECT total_messages, bot_count, non_bot_count FROM stats_summary WHERE id = 'global'`).WithContext(ctx).Consistency(gocql.One).Scan(&total, &bots, &nonBots); err != nil {
		log.Printf("error fetching summary: %v", err)
	}

	iter := cs.session.Query(`SELECT username FROM unique_users`).WithContext(ctx).Iter()
	totalUsers := 0
	for iter.Scan(new(string)) {
		totalUsers++
	}
	if err := iter.Close(); err != nil {
		log.Printf("error closing iterator for unique_users: %v", err)
	}

	serverCounts := make(map[string]int)
	iter = cs.session.Query(`SELECT server_url, count FROM server_counts`).WithContext(ctx).Iter()
	var url string
	var count int
	for iter.Scan(&url, &count) {
		serverCounts[url] = count
	}
	if err := iter.Close(); err != nil {
		log.Printf("error closing iterator for server_counts: %v", err)
	}

	return StatsSnapshot{
		Messages:      total,
		DistinctUsers: totalUsers,
		Bots:          bots,
		NonBots:       nonBots,
		ByServer:      serverCounts,
	}
}
