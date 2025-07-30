package skyline

import (
	"fmt"

	"github.com/gkoos/skyline/internal/algorithms"
	"github.com/gkoos/skyline/types"
)

var DNCConfig = types.DNCConfig{
	Threshold: 100,
	BatchSize: 100,
}

var SkyTreeConfig = types.SkyTreeConfig{
	// Add any SkyTree specific configuration here if needed
}

// Skyline computes the skyline from a static dataset using the specified algorithm.
// If algo is empty, defaults to "bnl".
func Skyline(points []types.Point, dims []string, prefs types.Preference, algo string) ([]types.Point, error) {
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
		result = algorithms.SkyTree(points, prefs, &SkyTreeConfig)
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
	return result, nil
}
