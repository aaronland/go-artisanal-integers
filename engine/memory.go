package engine

import (
	"context"
	"fmt"
	"sync/atomic"
)

type MemoryEngine struct {
	Engine
	increment int64
	offset    int64
	last      int64
}

func init() {

	ctx := context.Background()
	err := RegisterEngine(ctx, "memory", NewMemoryEngine)

	if err != nil {
		panic(err)
	}
}

func NewMemoryEngine(ctx context.Context, uri string) (Engine, error) {

	eng := &MemoryEngine{
		increment: 2,
		offset:    1,
		last:      0,
	}

	// PLEASE WRITE ME: check to see if we should read a value persisted to disk

	return eng, nil
}

func (eng *MemoryEngine) SetLastInt(i int64) error {

	last, err := eng.LastInt()

	if err != nil {
		return fmt.Errorf("Failed to retrieve last int, %w", err)
	}

	if last > i {
		return fmt.Errorf("%s is smaller than current last int", i)
	}

	atomic.StoreInt64(&eng.last, i)
	return nil
}

func (eng *MemoryEngine) SetKey(k string) error {
	return nil
}

func (eng *MemoryEngine) SetOffset(i int64) error {
	atomic.StoreInt64(&eng.offset, i)
	return nil
}

func (eng *MemoryEngine) SetIncrement(i int64) error {
	atomic.StoreInt64(&eng.increment, i)
	return nil
}

func (eng *MemoryEngine) NextInt() (int64, error) {
	next := atomic.AddInt64(&eng.last, eng.increment*eng.offset)
	return next, nil
}

func (eng *MemoryEngine) LastInt() (int64, error) {
	last := atomic.LoadInt64(&eng.last)
	return last, nil
}

func (eng *MemoryEngine) Close() error {
	return nil
}

// PLEASE WRITE ME

func (eng *MemoryEngine) persist(i int64) error {
	return fmt.Errorf("Not implemented.")
}
