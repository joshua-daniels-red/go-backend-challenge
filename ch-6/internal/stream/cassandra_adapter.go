package stream

import "github.com/gocql/gocql"

type CassandraSessionAdapter struct {
	sess *gocql.Session
}

func NewCassandraSessionAdapter(sess *gocql.Session) *CassandraSessionAdapter {
	return &CassandraSessionAdapter{sess: sess}
}

func (a *CassandraSessionAdapter) Query(stmt string, values ...interface{}) Query {
	return &CassandraQueryAdapter{q: a.sess.Query(stmt, values...)}
}

type CassandraQueryAdapter struct {
	q *gocql.Query
}

func (c *CassandraQueryAdapter) Exec() error {
	return c.q.Exec()
}

func (c *CassandraQueryAdapter) Iter() Iter {
	return &CassandraIterAdapter{i: c.q.Iter()}
}

type CassandraIterAdapter struct {
	i *gocql.Iter
}

func (c *CassandraIterAdapter) Scan(dest ...interface{}) bool {
	return c.i.Scan(dest...)
}

func (c *CassandraIterAdapter) Close() error {
	return c.i.Close()
}
