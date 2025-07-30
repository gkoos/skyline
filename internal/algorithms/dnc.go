package algorithms

import (
	"math/rand"
	"sort"
	"sync"

	"github.com/gkoos/skyline/internal/utilities"
	"github.com/gkoos/skyline/types"
)

var defaultDNCConfig = types.DNCConfig{Threshold: 100, BatchSize: 100}

func DivideAndConquer(data types.Dataset, prefs types.Preference, cfg *types.DNCConfig) types.Dataset {
	if cfg == nil {
		cfg = &defaultDNCConfig
	}

	// Apply BNL if small enough
	if len(data) <= cfg.Threshold {
		return BlockNestedLoop(data, prefs)
	}

	// Find dimension with largest range
	numDimensions := len(data[0])
	maxRange := 0.0
	splitDim := 0
	for d := range numDimensions {
		minVal, maxVal := data[0][d], data[0][d]
		for _, p := range data {
			if p[d] < minVal {
				minVal = p[d]
			}
			if p[d] > maxVal {
				maxVal = p[d]
			}
		}
		r := maxVal - minVal
		if r > maxRange {
			maxRange = r
			splitDim = d
		}
	}

	// Find median in splitDim (sort in place)
	sort.Slice(data, func(i, j int) bool {
		return data[i][splitDim] < data[j][splitDim]
	})
	medianIdx := len(data) / 2
	median := data[medianIdx][splitDim]

	// Partition points with random assignment for values equal to median
	var left, right types.Dataset
	for _, p := range data {
		if p[splitDim] < median {
			left = append(left, p)
		} else if p[splitDim] > median {
			right = append(right, p)
		} else {
			// p[splitDim] == median, randomly assign
			if rand.Intn(2) == 0 {
				left = append(left, p)
			} else {
				right = append(right, p)
			}
		}
	}

	// Parallelize recursive calls
	var leftSkyline, rightSkyline types.Dataset
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		leftSkyline = DivideAndConquer(left, prefs, cfg)
		wg.Done()
	}()
	go func() {
		rightSkyline = DivideAndConquer(right, prefs, cfg)
		wg.Done()
	}()
	wg.Wait()

	// Batch merge using cfg.BatchSize (symmetric merge)
	merged := make(types.Dataset, 0, len(leftSkyline)+len(rightSkyline))
	merged = appendNonDominated(merged, leftSkyline, rightSkyline, prefs, cfg.BatchSize)
	merged = appendNonDominated(merged, rightSkyline, leftSkyline, prefs, cfg.BatchSize)

	return merged
}

func appendNonDominated(merged types.Dataset, src, other types.Dataset, prefs types.Preference, batchSize int) types.Dataset {
	for i := 0; i < len(src); i += batchSize {
		end := i + batchSize
		if end > len(src) {
			end = len(src)
		}
		batch := src[i:end]
		for _, p := range batch {
			dominated := false
			for _, q := range other {
				if utilities.Dominates(q, p, prefs) {
					dominated = true
					break
				}
			}
			if !dominated {
				merged = append(merged, p)
			}
		}
	}
	return merged
}
