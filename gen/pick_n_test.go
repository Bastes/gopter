package gen

import (
	"slices"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestPickNProperties(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("picking zero items returns empty slice", prop.ForAll(
		func(items []int) bool {
			result := samplePickN(items, 0)
			// Should always return empty slice when picking 0 items
			return len(result) == 0
		},
		SliceOf(Int()),
	))

	properties.Property("picking from an empty slice returns empty slice", prop.ForAll(
		func(n int) bool {
			result := samplePickN([]int{}, n)
			// Should always return empty slice regardless of n
			return len(result) == 0
		},
		IntRange(0, 1000), // Generate random number of items to pick
	))

	properties.Property("picking more than available just returns all items", prop.ForAll(
		func(items []int, extra int) bool {
			// Ensure extra is at least 1
			if extra < 1 {
				extra = 1
			}
			n := len(items) + extra
			result := samplePickN(items, n)
			// Should return exactly the original slice (all items, nothing more)
			return slices.Equal(result, items)
		},
		SliceOf(Int()),
		IntRange(1, 100), // Generate a random number >= 1
	))

	properties.Property("picked items are subset of original", prop.ForAll(
		func(items []int) bool {
			if len(items) == 0 {
				return true
			}
			n := len(items) / 2
			result := samplePickN(items, n)

			for _, item := range result {
				if !slices.Contains(items, item) {
					return false
				}
			}
			return true
		},
		SliceOf(Int()),
	))

	properties.Property("picked items maintain original order", prop.ForAll(
		func() bool {
			items := []int{9, 3, 7, 1, 5, 8, 2, 6, 4}
			n := 4
			result := samplePickN(items, n)

			lastIndex := -1
			for _, item := range result {
				currentIndex := slices.Index(items, item)
				if currentIndex <= lastIndex {
					return false
				}
				lastIndex = currentIndex
			}
			return true
		},
	))

	properties.Property("result has exactly n elements", prop.ForAll(
		func(items []int) bool {
			if len(items) == 0 {
				return true
			}
			n := len(items) / 2
			if n == 0 {
				return true
			}
			result := samplePickN(items, n)

			return len(result) == n
		},
		SliceOf(IntRange(0, 100)),
	))

	properties.TestingRun(t)

	t.Run("eventually produces different selections", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		n := 5
		if !eventuallyProducesDifferentSelection(items, func() ([]int, bool) {
			picked, ok := PickN(items, n).Sample()
			if !ok || picked == nil {
				return nil, false
			}
			return picked.([]int), true
		}, 10) {
			t.Error("Expected generator to produce different selections, but got same selection every time")
		}
	})
}

// samplePickN is a test helper that samples from PickN and panics on failure.
// This simplifies property test code by eliminating repetitive error handling.
func samplePickN(items []int, n int) []int {
	picked, ok := PickN(items, n).Sample()
	if !ok || picked == nil {
		panic("PickN failed to generate sample")
	}
	return picked.([]int)
}

// eventuallyProducesDifferentSelection checks whether a generator function eventually produces
// a different selection from the original across multiple attempts.
func eventuallyProducesDifferentSelection[T comparable](original []T, generate func() ([]T, bool), attempts int) bool {
	for i := 0; i < attempts; i++ {
		result, ok := generate()
		if !ok {
			return false
		}

		if len(result) != len(original) {
			return true
		}

		for i := range original {
			if result[i] != original[i] {
				return true
			}
		}
	}
	return false
}
