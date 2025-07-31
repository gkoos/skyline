package algorithms

import (
	"testing"

	"github.com/gkoos/skyline/types"
)

func BenchmarkBNL_10000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset10000SmallSkyline4D, prefs)
	}
}

func BenchmarkDNC_10000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset10000SmallSkyline4D, prefs, nil)
	}
}

func BenchmarkSkyTree_10000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset10000SmallSkyline4D, prefs, nil)
	}
}

func BenchmarkBNL_5000OneDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset5000OneDominating, prefs)
	}
}

func BenchmarkDNC_5000OneDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset5000OneDominating, prefs, nil)
	}
}

func BenchmarkSkyTree_5000OneDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset5000OneDominating, prefs, nil)
	}
}

func BenchmarkBNL_5000CoupleDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset5000CoupleDominating, prefs)
	}
}

func BenchmarkDNC_5000CoupleDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset5000CoupleDominating, prefs, nil)
	}
}

func BenchmarkSkyTree_5000CoupleDominating(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset5000CoupleDominating, prefs, nil)
	}
}

func BenchmarkBNL_5000AllSame(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset5000AllSame, prefs)
	}
}

func BenchmarkDNC_5000AllSame(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset5000AllSame, prefs, nil)
	}
}

func BenchmarkSkyTree_5000AllSame(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset5000AllSame, prefs, nil)
	}
}

func BenchmarkBNL_1000OneDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset1000OneDominating4D, prefs)
	}
}

func BenchmarkDNC_1000OneDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset1000OneDominating4D, prefs, nil)
	}
}

func BenchmarkSkyTree_1000OneDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset1000OneDominating4D, prefs, nil)
	}
}

func BenchmarkBNL_1000CoupleDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset1000CoupleDominating4D, prefs)
	}
}

func BenchmarkDNC_1000CoupleDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset1000CoupleDominating4D, prefs, nil)
	}
}

func BenchmarkSkyTree_1000CoupleDominating4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset1000CoupleDominating4D, prefs, nil)
	}
}

func BenchmarkBNL_1000AllSame4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset1000AllSame4D, prefs)
	}
}

func BenchmarkDNC_1000AllSame4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset1000AllSame4D, prefs, nil)
	}
}

func BenchmarkSkyTree_1000AllSame4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset1000AllSame4D, prefs, nil)
	}
}

func BenchmarkBNL_100000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		BlockNestedLoop(Dataset100000SmallSkyline4D, prefs)
	}
}

func BenchmarkDNC_100000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset100000SmallSkyline4D, prefs, nil)
	}
}

func BenchmarkSkyTree_100000SmallSkyline4D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Max, types.Min, types.Max}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset100000SmallSkyline4D, prefs, nil)
	}
}

func BenchmarkDNC_200000ClusteredSmallSkyline8D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min}
	for i := 0; i < b.N; i++ {
		DivideAndConquer(Dataset200000ClusteredSmallSkyline8D, prefs, nil)
	}
}

func BenchmarkSkyTree_200000ClusteredSmallSkyline8D(b *testing.B) {
	prefs := types.Preference{types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min, types.Min}
	for i := 0; i < b.N; i++ {
		SkyTree(Dataset200000ClusteredSmallSkyline8D, prefs, nil)
	}
}
