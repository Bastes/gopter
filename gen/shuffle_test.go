package gen

import (
	"sort"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestShuffleProperties(t *testing.T) {
	t.Run("an empty slice stays empty", func(t *testing.T) {
		shuffled, ok := Shuffle([]int{}).Sample()
		if !ok {
			t.Fatal("Shuffle failed to generate sample for empty slice")
		}
		result := shuffled.([]int)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(result))
		}
	})

	properties := gopter.NewProperties(nil)

	properties.Property("a single element slice is left unchanged", prop.ForAll(
		func(n int) bool {
			shuffled, ok := Shuffle([]int{n}).Sample()
			if !ok || shuffled == nil {
				return false
			}
			result := shuffled.([]int)
			return len(result) == 1 && result[0] == n
		},
		Int(),
	))

	properties.Property("all elements of the original slice are given back", prop.ForAll(
		func(items []int) bool {
			shuffled, ok := Shuffle(items).Sample()
			if !ok || shuffled == nil {
				return false
			}
			result := shuffled.([]int)

			sortedOriginal := make([]int, len(items))
			copy(sortedOriginal, items)
			sort.Ints(sortedOriginal)

			sortedResult := make([]int, len(result))
			copy(sortedResult, result)
			sort.Ints(sortedResult)

			return intSlicesEqual(sortedOriginal, sortedResult)
		},
		SliceOf(Int()),
	))

	properties.TestingRun(t)

	t.Run("different ordering are (eventually) produced for multiple elements", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		if !eventuallyProducesDifferentOrder(items, func() ([]int, bool) {
			shuffled, ok := Shuffle(items).Sample()
			if !ok || shuffled == nil {
				return nil, false
			}
			return shuffled.([]int), true
		}, 10) {
			t.Error("Expected generator to produce different orderings, but got same ordering every time")
		}
	})
}

// eventuallyProducesDifferentOrder checks whether a generator function eventually produces
// an ordering different from the original across multiple attempts.
func eventuallyProducesDifferentOrder(original []int, generate func() ([]int, bool), attempts int) bool {
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
