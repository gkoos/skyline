package algorithms

import (
	"github.com/gkoos/skyline/types"
)

// Deterministic datasets for skyline algorithm tests

var (
	// 6400-point, 4D, 8 clusters, 8 outliers (all Min prefs)
	Dataset64Clusters4D = func() types.Dataset {
		data := make(types.Dataset, 6400)
		// 8 outlier points (skyline)
		outliers := []types.Point{
			{0, 10, 20, 30},
			{10, 20, 30, 0},
			{20, 30, 0, 10},
			{30, 0, 10, 20},
			{0, 20, 30, 10},
			{10, 30, 0, 20},
			{20, 0, 10, 30},
			{30, 10, 20, 0},
		}
		for c := 0; c < 8; c++ {
			// Set the outlier as the first point in each cluster
			data[c*800+0] = outliers[c]
			for i := 1; i < 800; i++ {
				data[c*800+i] = types.Point{
					outliers[c][0] + float64(i)/10.0,
					outliers[c][1] + float64(i)/10.0,
					outliers[c][2] + float64(i)/10.0,
					outliers[c][3] + float64(i)/10.0,
				}
			}
		}
		return data
	}()
	ExpectedSkyline64Clusters4D = func() types.Dataset {
		outliers := []types.Point{
			{0, 10, 20, 30},
			{10, 20, 30, 0},
			{20, 30, 0, 10},
			{30, 0, 10, 20},
			{0, 20, 30, 10},
			{10, 30, 0, 20},
			{20, 0, 10, 30},
			{30, 10, 20, 0},
		}
		return outliers
	}()

	// 10,000 points, 4D, small skyline
	Dataset10000SmallSkyline4D = func() types.Dataset {
		data := make(types.Dataset, 10000)
		// Most points are dominated, only a few are on the skyline
		for i := 0; i < 9990; i++ {
			// Clustered points, dominated by outliers
			data[i] = types.Point{float64(100 + i%10), float64(100 - i%10), float64(200 + i%10), float64(200 - i%10)}
		}
		// Add 10 outliers that will be the skyline
		for i := 0; i < 10; i++ {
			data[9990+i] = types.Point{float64(i), float64(1000 - i), float64(i), float64(2000 - i)}
		}
		return data
	}()

	ExpectedSkyline10000SmallSkyline4D = func() types.Dataset {
		data := make(types.Dataset, 10)
		for i := 0; i < 10; i++ {
			data[i] = types.Point{float64(i), float64(1000 - i), float64(i), float64(2000 - i)}
		}
		return data
	}()
	// 1000 points, one dominating (4D)

	Dataset1000OneDominating4D = func() types.Dataset {
		data := make(types.Dataset, 1000)
		// the first point, (0,0,0,0) dominates all others
		for i := 0; i < 1000; i++ {
			data[i] = types.Point{float64(i), float64(i), float64(i), float64(i)}
		}

		return data
	}()

	ExpectedSkyline1000OneDominating4D = types.Dataset{
		{0, 0, 0, 0},
	}

	// 1000 points, a couple dominating (4D)

	Dataset1000CoupleDominating4D = func() types.Dataset {
		data := make(types.Dataset, 1000)
		// 998 dominated points
		for i := 0; i < 998; i++ {
			data[i] = types.Point{float64(i + 1), float64(i + 1), float64(i + 1), float64(i + 1)}
		}
		// 2 dominating points
		data[998] = types.Point{0.5, 0.0, 0.5, 0.0}
		data[999] = types.Point{0.0, 0.5, 0.0, 0.5}
		return data
	}()

	ExpectedSkyline1000CoupleDominating4D = types.Dataset{
		{0.5, 0.0, 0.5, 0.0},
		{0.0, 0.5, 0.0, 0.5},
	}

	// 1000 points, all the same (4D)

	Dataset1000AllSame4D = func() types.Dataset {
		data := make(types.Dataset, 1000)
		for i := range data {
			data[i] = types.Point{7, 7, 7, 7}
		}
		return data
	}()

	ExpectedSkyline1000AllSame4D = func() types.Dataset {
		data := make(types.Dataset, 1000)
		for i := range data {
			data[i] = types.Point{7, 7, 7, 7}
		}
		return data
	}()
	// 5 elements, 2-3 dominating points

	Dataset5SomeDominating = types.Dataset{
		{1, 10}, // dominates all others
		{2, 9},  // dominated by {1,10}
		{3, 8},  // dominated by {1,10}
		{4, 7},  // dominated by {1,10}
		{5, 6},  // dominated by {1,10}
	}

	ExpectedSkyline5SomeDominating = types.Dataset{
		{1, 10},
	}

	// Edge cases

	DatasetEmpty = types.Dataset{}

	ExpectedSkylineEmpty = types.Dataset{}

	DatasetSingle = types.Dataset{
		{42, 42},
	}

	ExpectedSkylineSingle = types.Dataset{
		{42, 42},
	}

	DatasetAllSame = types.Dataset{
		{7, 7}, {7, 7}, {7, 7}, {7, 7}, {7, 7},
	}

	ExpectedSkylineAllSame = types.Dataset{
		{7, 7}, {7, 7}, {7, 7}, {7, 7}, {7, 7},
	}

	DatasetAllDominatedByOne = types.Dataset{
		{1, 10}, {2, 9}, {2, 8}, {3, 7}, {4, 6},
	}
	ExpectedSkylineAllDominatedByOne = types.Dataset{
		{1, 10},
	}

	// 5000 points, one dominating
	Dataset5000OneDominating = func() types.Dataset {
		data := make(types.Dataset, 5000)
		for i := 0; i < 4999; i++ {
			data[i] = types.Point{float64(100 + i), float64(100 - i)}
		}
		data[4999] = types.Point{0.0, 1000.0} // dominates all others
		return data
	}()
	ExpectedSkyline5000OneDominating = types.Dataset{
		{0, 1000},
	}

	// 5000 points, a couple dominating
	Dataset5000CoupleDominating = func() types.Dataset {
		data := make(types.Dataset, 5000)
		for i := 0; i < 4998; i++ {
			data[i] = types.Point{float64(100 + i), float64(100 - i)}
		}
		data[4998] = types.Point{0.0, 1000.0}
		data[4999] = types.Point{1.0, 999.0}
		return data
	}()
	ExpectedSkyline5000CoupleDominating = types.Dataset{
		{0, 1000},
	}

	// 5000 points, all the same
	Dataset5000AllSame = func() types.Dataset {
		data := make(types.Dataset, 5000)
		for i := range data {
			data[i] = types.Point{7, 7}
		}
		return data
	}()
	ExpectedSkyline5000AllSame = func() types.Dataset {
		data := make(types.Dataset, 5000)
		for i := range data {
			data[i] = types.Point{7, 7}
		}
		return data
	}()

	// 100,000 points, 4D, small skyline
	Dataset100000SmallSkyline4D = func() types.Dataset {
		data := make(types.Dataset, 100000)
		// Most points are dominated, only a few are on the skyline
		for i := 0; i < 99990; i++ {
			// Clustered points, dominated by outliers
			data[i] = types.Point{float64(100 + i%10), float64(100 - i%10), float64(200 + i%10), float64(200 - i%10)}
		}
		// Add 10 outliers that will be the skyline
		for i := 0; i < 10; i++ {
			data[99990+i] = types.Point{float64(i), float64(1000 - i), float64(i), float64(2000 - i)}
		}
		return data
	}()
	ExpectedSkyline100000SmallSkyline4D = func() types.Dataset {
		data := make(types.Dataset, 10)
		for i := 0; i < 10; i++ {
			data[i] = types.Point{float64(i), float64(1000 - i), float64(i), float64(2000 - i)}
		}
		return data
	}()

	// 200,000 points, 8D, clustered, small skyline
	Dataset200000ClusteredSmallSkyline8D = func() types.Dataset {
		data := make(types.Dataset, 200000)
		// 199,990 clustered points
		for i := 0; i < 199990; i++ {
			// Clustered around (100, 100, ..., 100) with small noise
			data[i] = types.Point{
				100 + float64(i%10), 100 + float64((i/10)%10), 100 + float64((i/100)%10), 100 + float64((i/1000)%10),
				100 + float64((i/10000)%10), 100 + float64((i/100000)%10), 100 + float64((i/1000000)%10), 100 + float64((i/10000000)%10),
			}
		}
		// 10 outliers that dominate all clustered points
		for i := 0; i < 10; i++ {
			data[199990+i] = types.Point{
				float64(i), float64(i), float64(i), float64(i), float64(i), float64(i), float64(i), float64(i),
			}
		}
		return data
	}()
	ExpectedSkyline200000ClusteredSmallSkyline8D = func() types.Dataset {
		data := make(types.Dataset, 10)
		for i := 0; i < 10; i++ {
			data[i] = types.Point{
				float64(i), float64(i), float64(i), float64(i), float64(i), float64(i), float64(i), float64(i),
			}
		}
		return data
	}()
	// 2000 points, 8D, 2 outliers, rest dominated
	Dataset2000SmallSkyline8D = func() types.Dataset {
		data := make(types.Dataset, 2000)
		// 1998 dominated points
		for i := 0; i < 1998; i++ {
			vals := make([]float64, 8)
			for d := 0; d < 8; d++ {
				vals[d] = 100 + float64(i%10) + float64(d)
			}
			data[i] = types.Point(vals)
		}
		// 2 outliers (skyline points)
		// These do not dominate each other, but both dominate all others
		data[1998] = types.Point{0, 0, 0, 0, 10, 10, 10, 10}
		data[1999] = types.Point{10, 10, 10, 10, 0, 0, 0, 0}
		return data
	}()
	ExpectedSkyline2000SmallSkyline8D = types.Dataset{
		{0, 0, 0, 0, 10, 10, 10, 10},
		{10, 10, 10, 10, 0, 0, 0, 0},
	}

	// 8 points, 8D, all points on the skyline and not equal (anti-chain)
	Dataset2000AllSkyline8D = func() types.Dataset {
		N := 8
		data := make(types.Dataset, N)
		for i := 0; i < N; i++ {
			data[i] = types.Point{
				float64(i), float64(i + 1), float64(i + 2), float64(i + 3),
				float64(N - i), float64(N - i - 1), float64(N - i - 2), float64(N - i - 3),
			}
		}
		return data
	}()
	ExpectedSkyline2000AllSkyline8D = Dataset2000AllSkyline8D

	// 2000 points, 8D, all points are equal
	Dataset2000AllEqual8D = func() types.Dataset {
		data := make(types.Dataset, 2000)
		for i := range data {
			data[i] = types.Point{7, 7, 7, 7, 7, 7, 7, 7}
		}
		return data
	}()
	ExpectedSkyline2000AllEqual8D = Dataset2000AllEqual8D
)
