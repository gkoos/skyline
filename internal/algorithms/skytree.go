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

// SkyTree computes the skyline using the SkyTree algorithm.
func SkyTree(data []types.Point, prefs types.Preference, cfg SkyTreeConfig) []types.Point {
	if len(data) <= cfg.BNLSwitchThreshold {
		return BNL(data, prefs, BNLConfig{Epsilon: cfg.Epsilon})
	}
	var pivots []types.Point
	for _, p := range data {
		dominated := false
		for _, q := range pivots {
			if utilities.DominatesEpsilon(q, p, prefs, cfg.Epsilon) {
				dominated = true
				break
			}
		}
		if !dominated {
			pivots = append(pivots, p)
		}
	}
	return blockNestedLoop(pivots, prefs)
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
