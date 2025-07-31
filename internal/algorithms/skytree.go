package algorithms

import (
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gkoos/skyline/internal/utilities"
	"github.com/gkoos/skyline/types"
)

// Semaphore for parallelization control
var maxParallelism int
var parallelSem chan struct{}

func init() {
	maxParallelism = runtime.NumCPU()
	parallelSem = make(chan struct{}, maxParallelism)
}

// dominanceKeyString creates a unique string key for a, b, and prefs
func dominanceKeyString(a, b types.Point, prefs types.Preference) string {
	var builder strings.Builder
	for _, pt := range [][]float64{a, b} {
		for _, v := range pt {
			builder.WriteString(strconv.FormatFloat(v, 'g', 17, 64))
			builder.WriteByte(',')
		}
	}
	for _, p := range prefs {
		builder.WriteString(strconv.Itoa(int(p)))
		builder.WriteByte(',')
	}
	return builder.String()
}

type dominanceCache struct {
	cache map[string]bool
	mu    sync.RWMutex
}

func newDominanceCache() *dominanceCache {
	return &dominanceCache{cache: make(map[string]bool)}
}

func (dc *dominanceCache) check(a, b types.Point, prefs types.Preference) (bool, bool) {
	key := dominanceKeyString(a, b, prefs)
	dc.mu.RLock()
	val, ok := dc.cache[key]
	dc.mu.RUnlock()
	return val, ok
}

func (dc *dominanceCache) store(a, b types.Point, prefs types.Preference, result bool) {
	key := dominanceKeyString(a, b, prefs)
	dc.mu.Lock()
	dc.cache[key] = result
	dc.mu.Unlock()
}

var defaultSkyTreeConfig = types.SkyTreeConfig{
	PivotSelector:     SelectMedianPivot,
	ParallelThreshold: 8,
}

// SkyTree computes the skyline using the SkyTree algorithm.
func SkyTree(data types.Dataset, prefs types.Preference, cfg *types.SkyTreeConfig) types.Dataset {
	if cfg == nil {
		cfg = &defaultSkyTreeConfig
	}
	skyline := skytreeRecWithDepth(data, prefs, cfg.PivotSelector, 0, cfg.MaxRecursionDepth)
	return skyline
}

// Recursive with depth limit and fallback
func skytreeRecWithDepth(data types.Dataset, prefs types.Preference, pivotSelector func(data types.Dataset, prefs types.Preference) types.Point, depth, maxDepth int) types.Dataset {
	n := len(data)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return data
	}
	if depth >= maxDepth {
		// Fallback to BlockNestedLoop
		return BlockNestedLoop(data, prefs)
	}
	pivot := pivotSelector(data, prefs)
	// Reuse buffer for remaining points
	var buf types.Dataset
	if cap(buf) < n {
		buf = make(types.Dataset, 0, n)
	} else {
		buf = buf[:0]
	}
	for _, pt := range data {
		if utilities.Dominates(pivot, pt, prefs) {
			continue
		}
		if utilities.Dominates(pt, pivot, prefs) {
			pivot = pt
			continue
		}
		buf = append(buf, pt)
	}
	remaining := buf
	partitions := make(map[int]types.Dataset)
	for _, pt := range remaining {
		mask := regionMaskBit(pt, pivot, prefs)
		partitions[mask] = append(partitions[mask], pt)
	}
	skyline := types.Dataset{pivot}
	partitionCount := len(partitions)
	var threshold int = 8
	parallel := partitionCount >= threshold
	dc := newDominanceCache()
	dominates := func(a, b types.Point, prefs types.Preference) bool {
		if val, ok := dc.check(a, b, prefs); ok {
			return val
		}
		res := utilities.Dominates(a, b, prefs)
		dc.store(a, b, prefs, res)
		return res
	}
	if parallel {
		var wg sync.WaitGroup
		var mu sync.Mutex
		results := make([]types.Dataset, 0, partitionCount)
		for _, subset := range partitions {
			if len(subset) == 0 || canPruneRegionCached(subset, skyline, prefs, dominates) {
				continue
			}
			parallelSem <- struct{}{} // Acquire slot
			wg.Add(1)
			go func(subset types.Dataset) {
				defer func() {
					<-parallelSem // Release slot
					wg.Done()
				}()
				pruned := pruneWithPivotCached(subset, pivot, prefs, dominates)
				childSky := skytreeRecWithDepth(pruned, prefs, pivotSelector, depth+1, maxDepth)
				mu.Lock()
				results = append(results, childSky)
				mu.Unlock()
			}(subset)
		}
		wg.Wait()
		for _, childSky := range results {
			skyline = mergeSkylineCached(skyline, childSky, prefs, dominates)
		}
	} else {
		for _, subset := range partitions {
			if len(subset) == 0 {
				continue
			}
			if canPruneRegionCached(subset, skyline, prefs, dominates) {
				continue
			}
			pruned := pruneWithPivotCached(subset, pivot, prefs, dominates)
			childSky := skytreeRecWithDepth(pruned, prefs, pivotSelector, depth+1, maxDepth)
			skyline = mergeSkylineCached(skyline, childSky, prefs, dominates)
		}
	}
	return skyline
}

// selectMedianPivot chooses the point closest to the median in all dimensions
func SelectMedianPivot(data types.Dataset, prefs types.Preference) types.Point {
	n := len(data)
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
	bestDist := euclideanDistance(data[0], medians)
	for _, pt := range data[1:] {
		dist := euclideanDistance(pt, medians)
		if dist < bestDist {
			best = pt
			bestDist = dist
		}
	}
	return best
}

func euclideanDistance(pt types.Point, center []float64) float64 {
	sum := 0.0
	for i := range pt {
		d := pt[i] - center[i]
		sum += d * d
	}
	return sum
}

// Cached versions of helpers

func canPruneRegionCached(region types.Dataset, skyline types.Dataset, prefs types.Preference, dominates func(a, b types.Point, prefs types.Preference) bool) bool {
	for _, pt := range region {
		dominated := false
		for _, s := range skyline {
			if dominates(s, pt, prefs) {
				dominated = true
				break
			}
		}
		if !dominated {
			return false // At least one point is not dominated, don't prune
		}
	}
	return true // All points are dominated, prune region
}

func pruneWithPivotCached(data types.Dataset, pivot types.Point, prefs types.Preference, dominates func(a, b types.Point, prefs types.Preference) bool) types.Dataset {
	// Reuse buffer for pruned points
	var buf types.Dataset
	if cap(buf) < len(data) {
		buf = make(types.Dataset, 0, len(data))
	} else {
		buf = buf[:0]
	}
	for _, pt := range data {
		if !dominates(pivot, pt, prefs) {
			buf = append(buf, pt)
		}
	}
	return buf
}

func mergeSkylineCached(a, b types.Dataset, prefs types.Preference, dominates func(a, b types.Point, prefs types.Preference) bool) types.Dataset {
	// Reuse buffer for result
	result := append(a[:0:0], a...)
	for _, pt := range b {
		dominated := false
		for _, s := range result {
			if dominates(s, pt, prefs) {
				dominated = true
				break
			}
		}
		if !dominated {
			// Remove any points in result dominated by pt
			buf := result[:0]
			for _, s := range result {
				if !dominates(pt, s, prefs) {
					buf = append(buf, s)
				}
			}
			buf = append(buf, pt)
			result = buf
		}
	}
	return result
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
