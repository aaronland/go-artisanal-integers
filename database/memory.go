package database

import (
	"context"
	"fmt"
	"sync/atomic"
)

type MemoryDatabase struct {
	Database
	increment int64
	offset    int64
	last      int64
}

func init() {

	ctx := context.Background()
	err := RegisterDatabase(ctx, "memory", NewMemoryDatabase)

	if err != nil {
		panic(err)
	}
}

func NewMemoryDatabase(ctx context.Context, uri string) (Database, error) {

	db := &MemoryDatabase{
		increment: 2,
		offset:    1,
		last:      0,
	}

	err := SetParametersFromURI(ctx, db, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to set parameters, %w", err)
	}

	return db, nil
}

func (db *MemoryDatabase) SetLastInt(ctx context.Context, i int64) error {

	last, err := db.LastInt(ctx)

	if err != nil {
		return fmt.Errorf("Failed to retrieve last int, %w", err)
	}

	if last > i {
		return fmt.Errorf("%d is smaller than current last int", i)
	}

	atomic.StoreInt64(&db.last, i)
	return nil
}

func (db *MemoryDatabase) SetOffset(ctx context.Context, i int64) error {
	atomic.StoreInt64(&db.offset, i)
	return nil
}

func (db *MemoryDatabase) SetIncrement(ctx context.Context, i int64) error {
	atomic.StoreInt64(&db.increment, i)
	return nil
}

func (db *MemoryDatabase) NextInt(ctx context.Context) (int64, error) {
	next := atomic.AddInt64(&db.last, db.increment*db.offset)
	return next, nil
}

func (db *MemoryDatabase) LastInt(ctx context.Context) (int64, error) {
	last := atomic.LoadInt64(&db.last)
	return last, nil
}

func (db *MemoryDatabase) Close(ctx context.Context) error {
	return nil
}
