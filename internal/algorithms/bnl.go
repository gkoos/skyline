package algorithms

import (
	"github.com/gkoos/skyline/internal/utilities"
	"github.com/gkoos/skyline/types"
)

type BNLConfig = types.BNLConfig

// BlockNestedLoop is a re-export for compatibility with static.go and tests.
func BlockNestedLoop(data []types.Point, prefs types.Preference) []types.Point {
	return BNL(data, prefs, BNLConfig{Epsilon: 0})
}

func BNL(data []types.Point, prefs types.Preference, cfg BNLConfig) []types.Point {
	var skyline []types.Point
	for _, p := range data {
		dominated := false
		for i := 0; i < len(skyline); {
			if utilities.DominatesEpsilon(skyline[i], p, prefs, cfg.Epsilon) {
				dominated = true
				break
			} else if utilities.DominatesEpsilon(p, skyline[i], prefs, cfg.Epsilon) {
				skyline = append(skyline[:i], skyline[i+1:]...)
			} else {
				i++
			}
		}
		if !dominated {
			skyline = append(skyline, p)
		}
	}
	return skyline
}
