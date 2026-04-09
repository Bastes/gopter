package gen

import "github.com/leanovate/gopter"

// Shuffle takes a slice and returns a generator that produces shuffled copies
func Shuffle[T any](items []T) gopter.Gen {
	n := len(items)
	if n == 0 {
		return Const([]T{})
	}

	gens := make([]gopter.Gen, n-1)
	for i := range n - 1 {
		gens[i] = IntRange(0, i+1)
	}

	return gopter.CombineGens(gens...).Map(func(swaps []any) []T {
		result := make([]T, n)
		copy(result, items)

		for i := n - 1; i > 0; i-- {
			j := swaps[i-1].(int)
			result[i], result[j] = result[j], result[i]
		}

		return result
	})
}
