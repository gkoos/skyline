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
	// Add fields as needed for SkyTree configuration
}
