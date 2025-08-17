package algorithms

import (
	"testing"

	"github.com/gkoos/skyline/types"
)

func TestSkyTree_Skyline_Large(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Dataset
		expected types.Dataset
		prefs    types.Preference
	}{
		{"5SomeDominating", Dataset5SomeDominating, ExpectedSkyline5SomeDominating, types.Preference{types.Min, types.Max}},
		{"Empty", DatasetEmpty, ExpectedSkylineEmpty, types.Preference{types.Min, types.Max}},
		{"Single", DatasetSingle, ExpectedSkylineSingle, types.Preference{types.Min, types.Max}},
		{"AllSame", DatasetAllSame, ExpectedSkylineAllSame, types.Preference{types.Min, types.Max}},
		{"AllDominatedByOne", DatasetAllDominatedByOne, ExpectedSkylineAllDominatedByOne, types.Preference{types.Min, types.Max}},
		{"5000OneDominating", Dataset5000OneDominating, ExpectedSkyline5000OneDominating, types.Preference{types.Min, types.Max}},
		{"5000CoupleDominating", Dataset5000CoupleDominating, ExpectedSkyline5000CoupleDominating, types.Preference{types.Min, types.Max}},
		{"5000AllSame", Dataset5000AllSame, ExpectedSkyline5000AllSame, types.Preference{types.Min, types.Max}},
		// 4D cases
		{"1000OneDominating4D", Dataset1000OneDominating4D, ExpectedSkyline1000OneDominating4D, types.Preference{types.Min, types.Min, types.Min, types.Min}},
		{"1000CoupleDominating4D", Dataset1000CoupleDominating4D, ExpectedSkyline1000CoupleDominating4D, types.Preference{types.Min, types.Min, types.Min, types.Min}},
		{"1000AllSame4D", Dataset1000AllSame4D, ExpectedSkyline1000AllSame4D, types.Preference{types.Min, types.Min, types.Min, types.Min}},
		// 8D cases
		{"2000SmallSkyline8D", Dataset2000SmallSkyline8D, ExpectedSkyline2000SmallSkyline8D, types.Preference{types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min}},
		{"2000AllSkyline8D", Dataset2000AllSkyline8D, ExpectedSkyline2000AllSkyline8D, types.Preference{types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min}},
		{"2000AllEqual8D", Dataset2000AllEqual8D, ExpectedSkyline2000AllEqual8D, types.Preference{types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min}},
		// Synthetic cluster/outlier test
		{"64Clusters4D", Dataset64Clusters4D, ExpectedSkyline64Clusters4D, types.Preference{types.Min, types.Min, types.Min, types.Min}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SkyTree(tc.input, tc.prefs, DefaultSkyTreeConfig)
			if !equalSkylineSet(result, tc.expected) {
				t.Errorf("SkyTree skyline incorrect for %s: got %v, want %v", tc.name, result, tc.expected)
			}
		})
	}
}
