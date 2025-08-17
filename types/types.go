package types

type Point []float64

type Dataset []Point

// Preference specifies per-dimension optimization (Min or Max)
type Preference []Order

type Order int

const (
	Min Order = iota
	Max
	Ignore // Skip this dimension in dominance comparisons
)

type DNCConfig struct {
	Threshold int
	BatchSize int
	Epsilon   float64 // Relaxed dominance tolerance
}

type SkyTreeConfig struct {
	PivotSelector      func(data Dataset, prefs Preference) Point
	ParallelThreshold  int     // Minimum number of partitions to parallelize
	MaxRecursionDepth  int     // Maximum allowed recursion depth for SkyTree
	BNLSwitchThreshold int     // Switch to BNL if len(data) <= this
	WorkerPoolSize     int     // Number of workers for parallel processing (0 = all available cores)
	Epsilon            float64 // Relaxed dominance tolerance
}

type BNLConfig struct {
	Epsilon float64 // Relaxed dominance tolerance
}
