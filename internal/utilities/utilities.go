package utilities

import "github.com/gkoos/skyline/types"

// Returns true if a dominates b according to the given preferences.
func Dominates(a, b types.Point, prefs types.Preference) bool {
	anyBetter := false

	for dim, order := range prefs {
		if order == types.Ignore {
			continue
		}

		av, bv := a[dim], b[dim]
		if order == types.Min {
			if av > bv {
				return false
			}
			if av < bv {
				anyBetter = true
			}
		} else { // types.Max
			if av < bv {
				return false
			}
			if av > bv {
				anyBetter = true
			}
		}
	}
	return anyBetter
}
