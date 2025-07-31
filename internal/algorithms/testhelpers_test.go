package algorithms

import (
	"reflect"

	"github.com/gkoos/skyline/types"
)

// equalSkylineSet compares two datasets as sets (order-insensitive)
func equalSkylineSet(a, b types.Dataset) bool {
	if len(a) != len(b) {
		return false
	}
	matched := make([]bool, len(b))
	for _, pa := range a {
		found := false
		for j, pb := range b {
			if matched[j] {
				continue
			}
			if reflect.DeepEqual(pa, pb) {
				matched[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
