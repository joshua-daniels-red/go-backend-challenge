package stream_test

import (
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-7/internal/stream"
	"github.com/stretchr/testify/assert"
)

type mockQuery struct {
	execFunc func() error
	iter     stream.Iter
}

func (m *mockQuery) Exec() error {
	if m.execFunc != nil {
		return m.execFunc()
	}
	return nil
}

func (m *mockQuery) Iter() stream.Iter {
	return m.iter
}

type mockIter struct {
	data     [][2]interface{}
	index    int
	closeErr error
}

func (m *mockIter) Scan(dest ...interface{}) bool {
	if m.index >= len(m.data) {
		return false
	}
	row := m.data[m.index]
	*dest[0].(*string) = row[0].(string)
	*dest[1].(*int) = row[1].(int)
	m.index++
	return true
}

func (m *mockIter) Close() error {
	return m.closeErr
}

type mockSession struct {
	calledQueries []string
	iter          stream.Iter
	queryOverride func(string, ...interface{}) stream.Query
}

func (m *mockSession) Query(stmt string, values ...interface{}) stream.Query {
	m.calledQueries = append(m.calledQueries, stmt)
	if m.queryOverride != nil {
		return m.queryOverride(stmt, values...)
	}
	return &mockQuery{iter: m.iter}
}

func TestCassandraStats_Record(t *testing.T) {
	mock := &mockSession{}
	stats := stream.NewCassandraStats(mock)

	event := stream.Event{Domain: "en.wikipedia.org", User: "alice"}
	stats.Record(event)

	assert.Len(t, mock.calledQueries, 2)
	assert.Contains(t, mock.calledQueries[0], "UPDATE stats_by_domain")
	assert.Contains(t, mock.calledQueries[1], "UPDATE stats_by_user")
}

func TestCassandraStats_GetSnapshot(t *testing.T) {
	iter := &mockIter{
		data: [][2]interface{}{
			{"en.wikipedia.org", 5},
			{"de.wikipedia.org", 3},
		},
	}

	mock := &mockSession{iter: iter}
	stats := stream.NewCassandraStats(mock)

	snapshot := stats.GetSnapshot()

	assert.Equal(t, 5, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 3, snapshot.ByDomain["de.wikipedia.org"])
}

func TestCassandraStats_Record_WithExecError(t *testing.T) {
	mock := &mockSession{
		queryOverride: func(stmt string, values ...interface{}) stream.Query {
			return &mockQuery{
				execFunc: func() error {
					return assert.AnError
				},
			}
		},
	}
	stats := stream.NewCassandraStats(mock)

	event := stream.Event{Domain: "test.com", User: "alice"}
	stats.Record(event)
}

func TestCassandraStats_GetSnapshot_WithCloseError(t *testing.T) {
	iter := &mockIter{
		data: [][2]interface{}{
			{"test.com", 1},
		},
		closeErr: assert.AnError, // this triggers the log error branch
	}

	mock := &mockSession{iter: iter}
	stats := stream.NewCassandraStats(mock)

	_ = stats.GetSnapshot() // capture log path on iter.Close()
}

func TestCassandraStats_GetSnapshot_FullCoverage(t *testing.T) {
	domainIter := &mockIter{
		data: [][2]interface{}{
			{"en.wikipedia.org", 3},
		},
	}

	userIter := &mockIter{
		data: [][2]interface{}{
			{"alice", 2},
		},
	}

	mock := &mockSession{
		queryOverride: func(stmt string, values ...interface{}) stream.Query {
			switch {
			case stmt == "SELECT domain, count FROM stats_by_domain":
				return &mockQuery{iter: domainIter}
			case stmt == "SELECT user, count FROM stats_by_user":
				return &mockQuery{iter: userIter}
			default:
				return &mockQuery{}
			}
		},
	}

	stats := stream.NewCassandraStats(mock)
	snapshot := stats.GetSnapshot()

	assert.Equal(t, 3, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 2, snapshot.ByUser["alice"]) // ✅ triggers the final uncovered line
}
