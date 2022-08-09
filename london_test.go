package artisanalinteger

import (
	"testing"
)

func TestIsLondonInteger(t *testing.T) {

	is_londonint := []int64{
		9,
		18,
	}

	not_londonint := []int64{
		2,
		3,
		28,
	}

	for _, i := range is_londonint {

		if !IsLondonInteger(i) {
			t.Fatalf("Expected %d to be London Integer", i)
		}
	}

	for _, i := range not_londonint {

		if IsLondonInteger(i) {
			t.Fatalf("Did not expect %d to be London Integer", i)
		}
	}
}
