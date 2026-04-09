package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// Shuffle takes a slice and returns a generator that produces shuffled copies
func Shuffle(items interface{}) gopter.Gen {
	itemsVal := reflect.ValueOf(items)
	if itemsVal.Kind() != reflect.Slice {
		panic("Shuffle requires a slice")
	}

	n := itemsVal.Len()
	if n == 0 {
		return Const(reflect.MakeSlice(itemsVal.Type(), 0, 0).Interface())
	}

	gens := make([]gopter.Gen, n-1)
	for i := 0; i < n-1; i++ {
		gens[i] = IntRange(0, i+1)
	}

	itemsCopy := items
	sliceType := itemsVal.Type()
	combinedGen := gopter.CombineGens(gens...)

	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		genResult := combinedGen(genParams)
		swapsValue, ok := genResult.Retrieve()
		if !ok {
			return gopter.NewGenResult(reflect.MakeSlice(sliceType, 0, 0).Interface(), gopter.NoShrinker)
		}

		swaps := swapsValue.([]interface{})
		itemsVal := reflect.ValueOf(itemsCopy)
		result := reflect.MakeSlice(sliceType, n, n)
		for i := 0; i < n; i++ {
			result.Index(i).Set(itemsVal.Index(i))
		}

		for i := n - 1; i > 0; i-- {
			j := swaps[i-1].(int)
			// Swap elements at positions i and j
			tmp := result.Index(i).Interface()
			result.Index(i).Set(result.Index(j))
			result.Index(j).Set(reflect.ValueOf(tmp))
		}

		return gopter.NewGenResult(result.Interface(), gopter.NoShrinker)
	}
}
