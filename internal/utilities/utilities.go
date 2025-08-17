package utilities

import "github.com/gkoos/skyline/types"

// DominatesEpsilon returns true if a dominates b according to the given preferences, allowing a tolerance epsilon.
func DominatesEpsilon(a, b types.Point, prefs types.Preference, epsilon float64) bool {
	anyBetter := false

	for dim, order := range prefs {
		if order == types.Ignore {
			continue
		}

		av, bv := a[dim], b[dim]
		if order == types.Min {
			if av > bv+epsilon {
				return false
			}
			if av < bv-epsilon {
				anyBetter = true
			}
		} else { // types.Max
			if av < bv-epsilon {
				return false
			}
			if av > bv+epsilon {
				anyBetter = true
			}
		}
	}
	return anyBetter
}
