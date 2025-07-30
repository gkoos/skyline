package utilities

import (
	"testing"

	"github.com/gkoos/skyline/types"
)

func TestDominates(t *testing.T) {
	cases := []struct {
		name     string
		a, b     types.Point
		prefs    types.Preference
		expected bool
	}{
		// Strict domination (all Min)
		{"StrictMinDomination", types.Point{1, 2}, types.Point{2, 3}, types.Preference{types.Min, types.Min}, true},

		// Strict domination (all Max)
		{"StrictMaxDomination", types.Point{3, 4}, types.Point{2, 3}, types.Preference{types.Max, types.Max}, true},

		// Equal points
		{"EqualPoints", types.Point{5, 5}, types.Point{5, 5}, types.Preference{types.Min, types.Min}, false},

		// Dominate in one, equal in other
		{"DominateOneEqualOther", types.Point{1, 2}, types.Point{2, 2}, types.Preference{types.Min, types.Min}, true},

		// Dominate in one, worse in other
		{"DominateOneWorseOther", types.Point{1, 4}, types.Point{2, 3}, types.Preference{types.Min, types.Min}, false},

		// Mixed Min/Max, domination
		{"MixedDomination", types.Point{1, 5}, types.Point{2, 3}, types.Preference{types.Min, types.Max}, true},

		// Mixed Min/Max, not domination
		{"MixedNotDomination", types.Point{1, 2}, types.Point{2, 3}, types.Preference{types.Min, types.Max}, false},

		// Negative values
		{"NegativeDomination", types.Point{-2, -3}, types.Point{-1, -2}, types.Preference{types.Min, types.Min}, true},

		// Zero values
		{"ZeroDomination", types.Point{0, 0}, types.Point{1, 1}, types.Preference{types.Min, types.Min}, true},

		// Floating point precision
		{"FloatPrecisionDomination", types.Point{1.0000001, 2.0}, types.Point{1.0000002, 2.0}, types.Preference{types.Min, types.Min}, true},

		// Dominate in some, not all
		{"DominateSomeNotAll", types.Point{1, 2, 4}, types.Point{2, 1, 3}, types.Preference{types.Min, types.Max, types.Min}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Dominates(tc.a, tc.b, tc.prefs)
			if result != tc.expected {
				t.Errorf("Dominates(%v, %v, %v) = %v, want %v", tc.a, tc.b, tc.prefs, result, tc.expected)
			}
		})
	}
}
