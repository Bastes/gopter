package gen_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/leanovate/gopter/gen"
)

func TestShuffle(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffleGen := gen.Shuffle(items)

	for i := 0; i < 100; i++ {
		value, ok := shuffleGen.Sample()
		if !ok {
			t.Error("Sample was not ok")
			continue
		}

		result, ok := value.([]int)
		if !ok {
			t.Errorf("Sample not slice of int: %#v", value)
			continue
		}

		if len(result) != len(items) {
			t.Errorf("Expected length %d, got %d: %v", len(items), len(result), result)
			continue
		}

		sortedOriginal := make([]int, len(items))
		copy(sortedOriginal, items)
		sort.Ints(sortedOriginal)

		sortedResult := make([]int, len(result))
		copy(sortedResult, result)
		sort.Ints(sortedResult)

		if !reflect.DeepEqual(sortedOriginal, sortedResult) {
			t.Errorf("Shuffled slice does not contain same elements as original.\nOriginal (sorted): %v\nResult (sorted): %v",
				sortedOriginal, sortedResult)
		}
	}
}

func TestShuffleEmpty(t *testing.T) {
	commonGeneratorTest(t, "shuffle empty", gen.Shuffle([]int{}), func(value interface{}) bool {
		result, ok := value.([]int)
		return ok && len(result) == 0
	})
}

func TestShuffleSingle(t *testing.T) {
	commonGeneratorTest(t, "shuffle single", gen.Shuffle([]int{42}), func(value interface{}) bool {
		result, ok := value.([]int)
		return ok && len(result) == 1 && result[0] == 42
	})
}

func TestShuffleVariety(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	gen := gen.Shuffle(items)

	first, ok := gen.Sample()
	if !ok {
		t.Fatal("First sample failed")
	}
	firstSlice := first.([]int)

	foundDifferent := false
	for i := 0; i < 10; i++ {
		sample, ok := gen.Sample()
		if !ok {
			t.Fatal("Sample failed")
		}
		slice := sample.([]int)
		if !reflect.DeepEqual(slice, firstSlice) {
			foundDifferent = true
			break
		}
	}

	if !foundDifferent {
		t.Error("Generator produced same ordering 10 times")
	}
}
