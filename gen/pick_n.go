package gen

import (
	"slices"

	"github.com/leanovate/gopter"
)

// PickN takes a slice and a number, and returns a generator that produces
// subsets of exactly n items in their original order, randomly selected without duplicates
func PickN[T any](items []T, number int) gopter.Gen {
	length := len(items)

	if number <= 0 || length == 0 {
		return Const([]T{})
	}

	if number >= length {
		return Const(items)
	}

	gens := make([]gopter.Gen, number)
	for i := 0; i < number; i++ {
		gens[i] = IntRange(i, length-1)
	}

	return gopter.CombineGens(gens...).Map(func(selections []any) []T {
		available := make([]int, length)
		for i := range available {
			available[i] = i
		}

		selected := make([]int, number)
		for i := 0; i < number; i++ {
			pickIdx := selections[i].(int)
			selected[i] = available[pickIdx]
			available[pickIdx] = available[i]
		}

		slices.Sort(selected)

		result := make([]T, number)
		for i, idx := range selected {
			result[i] = items[idx]
		}

		return result
	})
}
