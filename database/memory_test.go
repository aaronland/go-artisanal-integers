package database

import (
	"context"
	"testing"
)

func TestSimpleDatabase(t *testing.T) {

	ctx := context.Background()

	uri := "memory://"

	s, err := NewDatabase(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create database for %s, %v", uri, err)
	}

	offset := int64(2)
	increment := int64(2)
	last := int64(20)
	
	err = s.SetOffset(offset)

	if err != nil {
		t.Fatalf("Failed to set offset, %v", err)
	}

	err = s.SetIncrement(increment)

	if err != nil {
		t.Fatalf("Failed to set increment, %v", err)
	}

	err = s.SetLastInt(last)

	if err != nil {
		t.Fatalf("Failed to set last int, %v", err)
	}

	i, err := s.LastInt()

	if err != nil {
		t.Fatalf("Failed to get last int, %v", err)
	}

	if i != last {
		t.Fatalf("Unexpected last int: %d", i)
	}

	i, err = s.NextInt()

	if err != nil {
		t.Fatalf("Failed to get next int, %v", err)
	}

	if i != 24 {
		t.Fatalf("Unexpected next int: %d", i)
	}

}
