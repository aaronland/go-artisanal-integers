package service

import (
	"context"
	"testing"
)

func TestSimpleService(t *testing.T) {

	ctx := context.Background()

	uri := "simple://?engine=memory://"

	s, err := NewService(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create service for %s, %v", uri, err)
	}

	err = s.SetOffset(2)

	if err != nil {
		t.Fatalf("Failed to set offset, %v", err)
	}

	err = s.SetIncrement(2)

	if err != nil {
		t.Fatalf("Failed to set increment, %v", err)
	}

	err = s.SetLastInt(20)

	if err != nil {
		t.Fatalf("Failed to set last int, %v", err)
	}

	i, err := s.LastInt()

	if err != nil {
		t.Fatalf("Failed to get last int, %v", err)
	}

	if i != 20 {
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
