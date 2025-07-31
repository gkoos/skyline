package types

type Point []float64

type Dataset []Point

// Preference specifies per-dimension optimization (Min or Max)
type Preference []Order

type Order int

const (
	Min Order = iota
	Max
)

type DNCConfig struct {
	Threshold int
	BatchSize int
}

type SkyTreeConfig struct {
	PivotSelector     func(data Dataset, prefs Preference) Point
	ParallelThreshold int // Minimum number of partitions to parallelize
	MaxRecursionDepth int // Maximum allowed recursion depth for SkyTree
}
