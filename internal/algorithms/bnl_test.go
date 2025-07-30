package algorithms

import (
	"fmt"
	"testing"

	"github.com/gkoos/skyline/types"
)

func equalSkyline(a, b types.Dataset) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if pointKey(a[i]) != pointKey(b[i]) {
			return false
		}
	}
	return true
}

func pointKey(p types.Point) string {
	return fmt.Sprintf("%v", p)
}

func TestBNL_Skyline(t *testing.T) {
	prefs := types.Preference{types.Min, types.Max}
	tests := []struct {
		name     string
		input    types.Dataset
		expected types.Dataset
	}{
		{"5SomeDominating", Dataset5SomeDominating, ExpectedSkyline5SomeDominating},
		{"Empty", DatasetEmpty, ExpectedSkylineEmpty},
		{"Single", DatasetSingle, ExpectedSkylineSingle},
		{"AllSame", DatasetAllSame, ExpectedSkylineAllSame},
		{"AllDominatedByOne", DatasetAllDominatedByOne, ExpectedSkylineAllDominatedByOne},
		{"5000OneDominating", Dataset5000OneDominating, ExpectedSkyline5000OneDominating},
		{"5000CoupleDominating", Dataset5000CoupleDominating, ExpectedSkyline5000CoupleDominating},
		{"5000AllSame", Dataset5000AllSame, ExpectedSkyline5000AllSame},
		// 4D cases
		{"1000OneDominating4D", Dataset1000OneDominating4D, ExpectedSkyline1000OneDominating4D},
		{"1000CoupleDominating4D", Dataset1000CoupleDominating4D, ExpectedSkyline1000CoupleDominating4D},
		{"1000AllSame4D", Dataset1000AllSame4D, ExpectedSkyline1000AllSame4D},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := BlockNestedLoop(tc.input, prefs)
			if !equalSkyline(result, tc.expected) {
				t.Errorf("BNL skyline incorrect for %s: got %v, want %v", tc.name, result, tc.expected)
			}
		})
	}
}
