package algorithms

import "github.com/gkoos/skyline/types"

// Small deterministic datasets for skyline algorithm tests

var (
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
		for i := 0; i < 999; i++ {
			data[i] = types.Point{float64(100 + i), float64(100 - i), float64(200 + i), float64(200 - i)}
		}
		data[999] = types.Point{0.0, 1000.0, 0.0, 2000.0} // dominates all others
		return data
	}()
	ExpectedSkyline1000OneDominating4D = types.Dataset{
		{0, 1000, 0, 2000},
	}

	// 1000 points, a couple dominating (4D)
	Dataset1000CoupleDominating4D = func() types.Dataset {
		data := make(types.Dataset, 1000)
		for i := 0; i < 998; i++ {
			data[i] = types.Point{float64(100 + i), float64(100 - i), float64(200 + i), float64(200 - i)}
		}
		data[998] = types.Point{0.0, 1000.0, 0.0, 2000.0}
		data[999] = types.Point{1.0, 999.0, 1.0, 1999.0}
		return data
	}()
	ExpectedSkyline1000CoupleDominating4D = types.Dataset{
		{0, 1000, 0, 2000},
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
	DatasetEmpty         = types.Dataset{}
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
)
