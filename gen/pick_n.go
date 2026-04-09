package gen

import (
	"reflect"
	"sort"

	"github.com/leanovate/gopter"
)

// PickN takes a slice and a number, and returns a generator that produces
// subsets of exactly n items in their original order, randomly selected without duplicates
func PickN(items interface{}, number int) gopter.Gen {
	itemsVal := reflect.ValueOf(items)
	if itemsVal.Kind() != reflect.Slice {
		panic("PickN requires a slice")
	}

	length := itemsVal.Len()

	if number <= 0 || length == 0 {
		return Const(reflect.MakeSlice(itemsVal.Type(), 0, 0).Interface())
	}

	if number >= length {
		return Const(items)
	}

	gens := make([]gopter.Gen, number)
	for i := 0; i < number; i++ {
		gens[i] = IntRange(i, length-1)
	}

	itemsCopy := items
	sliceType := itemsVal.Type()
	combinedGen := gopter.CombineGens(gens...)

	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		genResult := combinedGen(genParams)
		selections, ok := genResult.Retrieve()
		if !ok {
			return gopter.NewGenResult(reflect.MakeSlice(sliceType, 0, 0).Interface(), gopter.NoShrinker)
		}

		selectionsSlice := selections.([]interface{})
		itemsVal := reflect.ValueOf(itemsCopy)
		available := make([]int, length)
		for i := range available {
			available[i] = i
		}

		selected := make([]int, number)
		for i := 0; i < number; i++ {
			pickIdx := selectionsSlice[i].(int)
			selected[i] = available[pickIdx]
			available[pickIdx] = available[i]
		}

		sort.Ints(selected)

		result := reflect.MakeSlice(sliceType, number, number)
		for i, idx := range selected {
			result.Index(i).Set(itemsVal.Index(idx))
		}

		return gopter.NewGenResult(result.Interface(), gopter.NoShrinker)
	}
}
