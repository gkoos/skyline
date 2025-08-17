package skyline

import (
	"fmt"

	"github.com/gkoos/skyline/internal/algorithms"
	"github.com/gkoos/skyline/types"
)

// DNCConfig controls the configuration for the Divide & Conquer skyline algorithm.
// Modifying this variable changes the behavior of the D&C algorithm globally.
var DNCConfig = types.DNCConfig{
	Threshold: 100,
	BatchSize: 100,
}

// SkyTreeConfig controls the configuration for the SkyTree skyline algorithm.
// Modifying this variable changes the behavior of the SkyTree algorithm globally.
var SkyTreeConfig = types.SkyTreeConfig{
	PivotSelector:      algorithms.SelectMedianPivot,
	MaxRecursionDepth:  500,
	ParallelThreshold:  4,
	BNLSwitchThreshold: 1024,
	WorkerPoolSize:     0,
}

// Skyline computes the skyline from a static dataset using the specified algorithm.
// If algo is empty, defaults to "bnl".
func Skyline(points []types.Point, _ []string, prefs types.Preference, algo string) ([]types.Point, error) {
	if algo == "" {
		algo = "bnl"
	}

	var result []types.Point
	switch algo {
	case "bnl":
		result = algorithms.BlockNestedLoop(points, prefs)
	case "dnc":
		result = algorithms.DivideAndConquer(points, prefs, &DNCConfig)
	case "skytree":
		result = algorithms.SkyTree(points, prefs, SkyTreeConfig)
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
	return result, nil
}
