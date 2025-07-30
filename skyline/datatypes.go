package skyline

import (
	"github.com/gkoos/skyline/types"
)

// Point represents a multi-dimensional data point for skyline queries.
type Point = types.Point

// Dataset is a collection of Points used as input for skyline algorithms.
type Dataset = types.Dataset

// Preference maps each dimension to an optimization order (Min or Max).
type Preference = types.Preference

// Order specifies whether a dimension should be minimized or maximized.
type Order = types.Order

const (
	Min = types.Min // Minimize this dimension
	Max = types.Max // Maximize this dimension
)
