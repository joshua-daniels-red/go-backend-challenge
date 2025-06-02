package mocks

import(
	"errors"

 	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/stream"
)

type MockQuery struct {
	ExecFunc func() error
	iter     stream.Iter
}

func (m *MockQuery) Exec() error {
	if m.ExecFunc != nil {
		return m.ExecFunc()
	}
	return nil
}

func (m *MockQuery) Iter() stream.Iter {
	return m.iter
}

type MockIter struct {
	ScanFunc  func(dest ...interface{}) bool
	CloseFunc func() error
}

func (m *MockIter) Scan(dest ...interface{}) bool {
	return m.ScanFunc(dest...)
}

func (m *MockIter) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type mockSession struct {
	queryErr     bool
	iter         stream.Iter
	calledQueries []string
}

func (m *mockSession) Query(stmt string, values ...interface{}) stream.Query {
	m.calledQueries = append(m.calledQueries, stmt)
	return &MockQuery{
		ExecFunc: func() error {
			if m.queryErr {
				return errors.New("mock exec error")
			}
			return nil
		},
		iter: m.iter,
	}
}

func (m *mockSession) Close() {}
