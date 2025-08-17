package utilities

import (
	"testing"

	"github.com/gkoos/skyline/types"
)

func TestDominatesEpsilon_Zero(t *testing.T) {
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

		// Ignore dimension
		{"IgnoreDimension", types.Point{1, 2, 3}, types.Point{2, 1, 3}, types.Preference{types.Min, types.Ignore, types.Min}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := DominatesEpsilon(tc.a, tc.b, tc.prefs, 0)
			if result != tc.expected {
				t.Errorf("DominatesEpsilon(%v, %v, %v, 0) = %v, want %v", tc.a, tc.b, tc.prefs, result, tc.expected)
			}
		})
	}
}

// Additional tests: same dataset, different epsilons
// Placed at end of file per user request
func TestDominatesEpsilon_VaryingEpsilon(t *testing.T) {
	a := types.Point{1.0, 2.0}
	b := types.Point{1.01, 2.0}
	prefs := types.Preference{types.Min, types.Min}

	cases := []struct {
		name     string
		epsilon  float64
		expected bool
	}{
		{"EpsilonZero_Domination", 0, true},
		{"EpsilonSmall_Domination", 0.005, true},
		{"EpsilonCoversDiff_Domination", 0.02, false},
		{"EpsilonExactDiff_Domination", 0.01, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := DominatesEpsilon(a, b, prefs, tc.epsilon)
			if result != tc.expected {
				t.Errorf("DominatesEpsilon(%v, %v, %v, %v) = %v, want %v", a, b, prefs, tc.epsilon, result, tc.expected)
			}
		})
	}

	// Also test with Max preference
	a = types.Point{2.0, 5.0}
	b = types.Point{1.99, 5.0}
	prefs = types.Preference{types.Max, types.Max}
	cases = []struct {
		name     string
		epsilon  float64
		expected bool
	}{
		{"EpsilonZero_Domination_Max", 0, true},
		{"EpsilonSmall_Domination_Max", 0.005, true},
		{"EpsilonCoversDiff_Domination_Max", 0.02, false},
		{"EpsilonExactDiff_Domination_Max", 0.01, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := DominatesEpsilon(a, b, prefs, tc.epsilon)
			if result != tc.expected {
				t.Errorf("DominatesEpsilon(%v, %v, %v, %v) = %v, want %v", a, b, prefs, tc.epsilon, result, tc.expected)
			}
		})
	}
}
