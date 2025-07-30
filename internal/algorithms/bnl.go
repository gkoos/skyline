package algorithms

import (
	"github.com/gkoos/skyline/internal/utilities"
	"github.com/gkoos/skyline/types"
)

func BlockNestedLoop(data types.Dataset, prefs types.Preference) types.Dataset {
	n := len(data)
	if n == 0 {
		return nil
	}
	window := make([]bool, n) // true if candidate for skyline
	for i := range window {
		window[i] = true
	}
	for i := 0; i < n; i++ {
		if !window[i] {
			continue
		}
		for j := 0; j < n; j++ {
			if i == j || !window[j] {
				continue
			}
			if utilities.Dominates(data[j], data[i], prefs) {
				window[i] = false
				break
			}
		}
	}
	var result types.Dataset
	for i, ok := range window {
		if ok {
			result = append(result, data[i])
		}
	}
	return result
}
