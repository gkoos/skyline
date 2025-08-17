package algorithms

import (
	"sort"

	"github.com/gkoos/skyline/internal/utilities"
	"github.com/gkoos/skyline/types"
)

// DefaultSkyTreeConfig provides a default config for tests and static.go
var DefaultSkyTreeConfig = types.SkyTreeConfig{
	PivotSelector:      SelectMedianPivot,
	MaxRecursionDepth:  500,
	ParallelThreshold:  4,
	BNLSwitchThreshold: 1024,
	WorkerPoolSize:     0,
}

// SelectMedianPivot is a classic median pivot selector for static.go and tests
func SelectMedianPivot(data types.Dataset, _ types.Preference) types.Point {
	n := len(data)
	if n == 0 {
		return nil
	}
	dim := len(data[0])
	medians := make([]float64, dim)
	for i := 0; i < dim; i++ {
		vals := make([]float64, n)
		for j, pt := range data {
			vals[j] = pt[i]
		}
		sort.Float64s(vals)
		if n%2 == 0 {
			medians[i] = (vals[n/2-1] + vals[n/2]) / 2
		} else {
			medians[i] = vals[n/2]
		}
	}
	best := data[0]
	bestDist := 0.0
	for i := range best {
		d := best[i] - medians[i]
		bestDist += d * d
	}
	for _, pt := range data[1:] {
		dist := 0.0
		for i := range pt {
			d := pt[i] - medians[i]
			dist += d * d
		}
		if dist < bestDist {
			best = pt
			bestDist = dist
		}
	}
	return best
}

type SkyTreeConfig = types.SkyTreeConfig

// SkyTree computes the skyline using the SkyTree algorithm with pivot selection.
func SkyTree(data []types.Point, prefs types.Preference, cfg SkyTreeConfig) []types.Point {
	// Base cases
	n := len(data)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return data
	}
	if n <= cfg.BNLSwitchThreshold {
		return BNL(data, prefs, BNLConfig{Epsilon: cfg.Epsilon})
	}

	// Select pivot using the configured selector
	pivot := cfg.PivotSelector(data, prefs)
	if pivot == nil {
		return nil
	}

	// Partition data into points equal to pivot and the rest
	equalToPivot := make([]types.Point, 0, n)
	remaining := make([]types.Point, 0, n)
	for _, pt := range data {
		if isPointEqual(pt, pivot) {
			equalToPivot = append(equalToPivot, pt)
		} else {
			remaining = append(remaining, pt)
		}
	}

	// Partition remaining points by region relative to pivot
	partitions := make(map[int][]types.Point)
	for _, pt := range remaining {
		mask := regionMaskBit(pt, pivot, prefs)
		partitions[mask] = append(partitions[mask], pt)
	}

	// Recursively compute skylines for each partition
	var merged [][]types.Point
	for _, subset := range partitions {
		if len(subset) == 0 {
			continue
		}
		childSky := SkyTree(subset, prefs, cfg)
		merged = append(merged, childSky)
	}

	// Merge all child skylines and points equal to pivot
	result := make([]types.Point, 0, n)
	for _, sky := range merged {
		result = append(result, sky...)
	}
	result = append(result, equalToPivot...)

	// Final local skyline
	return blockNestedLoop(result, prefs)
}

// isPointEqual checks if two points are exactly equal
func isPointEqual(a, b types.Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// regionMaskBit encodes the region of pt relative to pivot as an integer bitmask
func regionMaskBit(pt, pivot types.Point, prefs types.Preference) int {
	mask := 0
	for i := range pt {
		if pt[i] == pivot[i] {
			continue // bit stays 0
		}
		if (prefs[i] == types.Min && pt[i] < pivot[i]) || (prefs[i] == types.Max && pt[i] > pivot[i]) {
			mask |= 1 << i // set bit i
		}
		// else: bit stays 0
	}
	return mask
}

// blockNestedLoop is a local version for internal use
func blockNestedLoop(data []types.Point, prefs types.Preference) []types.Point {
	n := len(data)
	if n == 0 {
		return nil
	}
	window := make([]bool, n)
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
			if utilities.DominatesEpsilon(data[j], data[i], prefs, 0) {
				window[i] = false
				break
			}
		}
	}
	var result []types.Point
	for i, ok := range window {
		if ok {
			result = append(result, data[i])
		}
	}
	return result
}
