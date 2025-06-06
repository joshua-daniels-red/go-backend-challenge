package stream

import (
	"context"
	"sync"
	"time"
)

type Batcher struct {
	store         StatsStore
	batchSize     int
	flushInterval time.Duration

	mu      sync.Mutex
	buffer  []Event
	ticker  *time.Ticker
	wg      sync.WaitGroup
	flushCh chan struct{}
}

func NewBatcher(store StatsStore, batchSize int, flushInterval time.Duration) *Batcher {
	b := &Batcher{
		store:         store,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		buffer:        make([]Event, 0, batchSize),
		ticker:        time.NewTicker(flushInterval),
		flushCh:       make(chan struct{}, 1),
	}
	return b
}

func (b *Batcher) Start(ctx context.Context) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-ctx.Done():
				b.flush()
				return
			case <-b.ticker.C:
				b.flush()
			case <-b.flushCh:
				b.flush()
			}
		}
	}()
}

func (b *Batcher) Add(event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer = append(b.buffer, event)

	if len(b.buffer) >= b.batchSize {
		select {
		case b.flushCh <- struct{}{}:
		default:
		}
	}
}

func (b *Batcher) flush() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buffer) == 0 {
		return nil
	}

	toFlush := make([]Event, len(b.buffer))
	copy(toFlush, b.buffer)

	b.store.RecordMany(toFlush)
	b.buffer = b.buffer[:0]

	return toFlush
}

func (b *Batcher) Stop() {
	b.ticker.Stop()
	close(b.flushCh)
	b.wg.Wait()
}

func (b *Batcher) FlushIfThresholdMet() []Event {
	b.mu.Lock()
	shouldFlush := len(b.buffer) >= b.batchSize
	b.mu.Unlock()

	if shouldFlush {
		return b.flush()
	}
	return nil
}
