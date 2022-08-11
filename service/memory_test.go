package service

import (
	"context"
	"testing"
)

func TestSimpleService(t *testing.T) {

	ctx := context.Background()

	uri := "memory://"

	s, err := NewService(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create service for %s, %v", uri, err)
	}

	offset := int64(2)
	increment := int64(2)
	last := int64(20)

	err = s.SetOffset(ctx, offset)

	if err != nil {
		t.Fatalf("Failed to set offset, %v", err)
	}

	err = s.SetIncrement(ctx, increment)

	if err != nil {
		t.Fatalf("Failed to set increment, %v", err)
	}

	err = s.SetLastInt(ctx, last)

	if err != nil {
		t.Fatalf("Failed to set last int, %v", err)
	}

	i, err := s.LastInt(ctx)

	if err != nil {
		t.Fatalf("Failed to get last int, %v", err)
	}

	if i != last {
		t.Fatalf("Unexpected last int: %d", i)
	}

	i, err = s.NextInt(ctx)

	if err != nil {
		t.Fatalf("Failed to get next int, %v", err)
	}

	if i != 24 {
		t.Fatalf("Unexpected next int: %d", i)
	}

}
