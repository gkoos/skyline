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
// Worker pool for parallel recursion
type workerPool struct {
	sem chan struct{}
}

func newWorkerPool(size int) *workerPool {
	if size <= 0 {
		size = runtime.NumCPU()
	}
	return &workerPool{sem: make(chan struct{}, size)}
}

func (wp *workerPool) acquire() {
	wp.sem <- struct{}{}
}

func (wp *workerPool) release() {
	<-wp.sem
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
	PivotSelector:      SelectMedianPivot,
	MaxRecursionDepth:  500,
	ParallelThreshold:  4,
	BNLSwitchThreshold: 1024, // Default: switch to BNL if <= 1024 points
	WorkerPoolSize:     0,    // Default: use all available CPU cores
}

// SkyTree computes the skyline using the SkyTree algorithm.
func SkyTree(data types.Dataset, prefs types.Preference, cfg *types.SkyTreeConfig) types.Dataset {
	if cfg == nil {
		cfg = &defaultSkyTreeConfig
	}
	pool := newWorkerPool(cfg.WorkerPoolSize)
	skyline := skytreeRecWithDepthPool(data, prefs, cfg.PivotSelector, 0, cfg.MaxRecursionDepth, pool, cfg.BNLSwitchThreshold)
	return skyline
}

// Recursive with worker pool and deeper parallel recursion
func skytreeRecWithDepthPool(data types.Dataset, prefs types.Preference, pivotSelector func(data types.Dataset, prefs types.Preference) types.Point, depth, maxDepth int, pool *workerPool, bnlSwitchThreshold int) types.Dataset {
	n := len(data)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return data
	}
	if n <= bnlSwitchThreshold {
		return BlockNestedLoop(data, prefs)
	}
	// Check for anti-chain: all points are mutually non-dominating
	antiChain := true
	for i := 0; i < n && antiChain; i++ {
		for j := i + 1; j < n; j++ {
			if utilities.Dominates(data[i], data[j], prefs) || utilities.Dominates(data[j], data[i], prefs) {
				antiChain = false
				break
			}
		}
	}
	if antiChain {
		return data
	}
	if depth >= maxDepth {
		// Fallback to BlockNestedLoop
		return BlockNestedLoop(data, prefs)
	}
	pivot := pivotSelector(data, prefs)
	// Partition points and collect those equal to the pivot
	equalToPivot := make(types.Dataset, 0, n)
	remaining := make(types.Dataset, 0, n)
	for _, pt := range data {
		equal := true
		if len(pt) == len(pivot) {
			for i := range pt {
				if pt[i] != pivot[i] {
					equal = false
					break
				}
			}
		} else {
			equal = false
		}
		if equal {
			equalToPivot = append(equalToPivot, pt)
		} else {
			remaining = append(remaining, pt)
		}
	}
	partitions := make(map[int]types.Dataset)
	for _, pt := range remaining {
		mask := regionMaskBit(pt, pivot, prefs)
		partitions[mask] = append(partitions[mask], pt)
	}

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
	// Gather all child skylines
	childSkylines := make([]types.Dataset, partitionCount)
	keys := make([]int, 0, partitionCount)
	for k := range partitions {
		keys = append(keys, k)
	}
	var wg sync.WaitGroup
	if parallel {
		for idx, k := range keys {
			subset := partitions[k]
			if len(subset) == 0 {
				continue
			}
			pruned := pruneWithPivotCached(subset, pivot, prefs, dominates)
			wg.Add(1)
			pool.acquire()
			go func(i int, prunedSubset types.Dataset) {
				defer func() {
					pool.release()
					wg.Done()
				}()
				childSky := skytreeRecWithDepthPool(prunedSubset, prefs, pivotSelector, depth+1, maxDepth, pool, bnlSwitchThreshold)
				childSkylines[i] = childSky
			}(idx, pruned)
		}
		wg.Wait()
	} else {
		for idx, k := range keys {
			subset := partitions[k]
			if len(subset) == 0 {
				continue
			}
			pruned := pruneWithPivotCached(subset, pivot, prefs, dominates)
			childSky := skytreeRecWithDepthPool(pruned, prefs, pivotSelector, depth+1, maxDepth, pool, bnlSwitchThreshold)
			childSkylines[idx] = childSky
		}
	}
	// Parallel merge of child skylines using the worker pool
	merged := parallelMergeSkylines(childSkylines, prefs, pool)
	// Add points equal to the pivot
	merged = append(merged, equalToPivot...)
	// Compute the local skyline of the union
	localSkyline := BlockNestedLoop(merged, prefs)
	return localSkyline
}

// parallelMergeSkylines merges a slice of skylines in parallel using the worker pool
func parallelMergeSkylines(skylines []types.Dataset, prefs types.Preference, pool *workerPool) types.Dataset {
	if len(skylines) == 0 {
		return nil
	}
	if len(skylines) == 1 {
		return skylines[0]
	}
	// Iteratively merge in log2(N) stages
	curr := skylines
	for len(curr) > 1 {
		var wg sync.WaitGroup
		next := make([]types.Dataset, (len(curr)+1)/2)
		for i := 0; i < len(curr)/2; i++ {
			a, b := curr[2*i], curr[2*i+1]
			wg.Add(1)
			pool.acquire()
			go func(idx int, left, right types.Dataset) {
				defer func() {
					pool.release()
					wg.Done()
				}()
				// Merge two skylines and compute their local skyline
				merged := make(types.Dataset, 0, len(left)+len(right))
				merged = append(merged, left...)
				merged = append(merged, right...)
				next[idx] = BlockNestedLoop(merged, prefs)
			}(i, a, b)
		}
		// If odd, carry the last one
		if len(curr)%2 == 1 {
			next[len(next)-1] = curr[len(curr)-1]
		}
		wg.Wait()
		curr = next
	}
	return curr[0]
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
