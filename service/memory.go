package service

import (
	"context"
	"fmt"
	"sync/atomic"
)

type MemoryService struct {
	Service
	increment int64
	offset    int64
	last      int64
}

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "memory", NewMemoryService)

	if err != nil {
		panic(err)
	}
}

func NewMemoryService(ctx context.Context, uri string) (Service, error) {

	s := &MemoryService{
		increment: 2,
		offset:    1,
		last:      0,
	}

	err := SetParametersFromURI(ctx, s, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to set parameters, %w", err)
	}

	return s, nil
}

func (s *MemoryService) SetLastInt(ctx context.Context, i int64) error {

	last, err := s.LastInt(ctx)

	if err != nil {
		return fmt.Errorf("Failed to retrieve last int, %w", err)
	}

	if last > i {
		return fmt.Errorf("%d is smaller than current last int", i)
	}

	atomic.StoreInt64(&s.last, i)
	return nil
}

func (s *MemoryService) SetOffset(ctx context.Context, i int64) error {
	atomic.StoreInt64(&s.offset, i)
	return nil
}

func (s *MemoryService) SetIncrement(ctx context.Context, i int64) error {
	atomic.StoreInt64(&s.increment, i)
	return nil
}

func (s *MemoryService) NextInt(ctx context.Context) (int64, error) {
	next := atomic.AddInt64(&s.last, s.increment*s.offset)
	return next, nil
}

func (s *MemoryService) LastInt(ctx context.Context) (int64, error) {
	last := atomic.LoadInt64(&s.last)
	return last, nil
}

func (s *MemoryService) Close(ctx context.Context) error {
	return nil
}
