// ...existing code...
package skyline

import (
	"github.com/gkoos/skyline/internal/utilities"
)

// Engine is the interface for dynamic skyline operations.
// Engine supports insert, update, delete, and retrieval of the current skyline set.
type Engine interface {
	Insert(Point)
	Update(Point, Point)
	Delete(Point)
	Skyline() []Point
}

// internal engine struct, all fields private
type engine struct {
	points  []Point
	dims    []string
	prefs   Preference
	algo    string
	skyline []Point // always up-to-date skyline set
}

// DynamicSkyline creates a new dynamic skyline Engine and calculates the initial skyline.
// DynamicSkyline returns an Engine that supports incremental skyline updates.
func DynamicSkyline(points []Point, dims []string, prefs Preference, algo string) (Engine, error) {
	e := &engine{
		points: points,
		dims:   dims,
		prefs:  prefs,
		algo:   algo,
	}
	// Compute initial skyline using the selected algorithm and current config
	result, err := Skyline(points, dims, prefs, algo)
	if err != nil {
		return nil, err
	}
	e.skyline = result
	return e, nil
}

// DynamicSkylineRaw creates a new dynamic skyline Engine using the provided points as the initial set, skipping skyline computation.
// If algo is empty, it defaults to "bnl" for later batch operations. This is useful for batch insertion or when the dataset is already known to be the skyline.
func DynamicSkylineRaw(points []Point, dims []string, prefs Preference, algo string) Engine {
	if algo == "" {
		algo = "bnl"
	}
	return &engine{
		points:  points,
		dims:    dims,
		prefs:   prefs,
		algo:    algo,
		skyline: append([]Point(nil), points...),
	}
}

// Insert adds a new point and updates the skyline incrementally.
func (e *engine) Insert(p Point) {
	e.points = append(e.points, p)

	// Optimized BNL: update skyline incrementally
	dominated := false
	var newSkyline []Point

	// Check if new point is dominated by any current skyline point
	for _, s := range e.skyline {
		if utilities.Dominates(s, p, e.prefs) {
			dominated = true
			break
		}
	}

	// If new point is dominated, skyline unchanged
	if dominated {
		return
	}

	// New point is not dominated, add it to skyline and remove any skyline points it dominates
	for _, s := range e.skyline {
		if utilities.Dominates(p, s, e.prefs) {
			// p dominates s, so s is not in new skyline
			continue
		}
		newSkyline = append(newSkyline, s)
	}
	newSkyline = append(newSkyline, p)
	e.skyline = newSkyline
}

// Update replaces an old point with a new one and updates the skyline.
func (e *engine) Update(old, new Point) {
	// Remove the old point from the dataset
	var updatedPoints []Point
	for _, pt := range e.points {
		if !equalPoint(pt, old) {
			updatedPoints = append(updatedPoints, pt)
		}
	}
	e.points = updatedPoints

	// Remove the old point from the skyline if present
	var updatedSkyline []Point
	for _, s := range e.skyline {
		if !equalPoint(s, old) {
			updatedSkyline = append(updatedSkyline, s)
		}
	}
	e.skyline = updatedSkyline

	// Insert the new point using the optimized BNL logic
	e.Insert(new)
}

// Delete removes a point and updates the skyline.
func (e *engine) Delete(p Point) {
	// Remove the point from the dataset
	var updatedPoints []Point
	for _, pt := range e.points {
		if !equalPoint(pt, p) {
			updatedPoints = append(updatedPoints, pt)
		}
	}
	e.points = updatedPoints

	// Remove the point from the skyline if present
	var updatedSkyline []Point
	for _, s := range e.skyline {
		if !equalPoint(s, p) {
			updatedSkyline = append(updatedSkyline, s)
		}
	}

	// For each point not in the skyline, check if it should now be added
	for _, candidate := range e.points {
		// If candidate is already in updatedSkyline, skip
		found := false
		for _, s := range updatedSkyline {
			if equalPoint(s, candidate) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		// Check if candidate is dominated by any skyline point
		dominated := false
		for _, s := range updatedSkyline {
			if utilities.Dominates(s, candidate, e.prefs) {
				dominated = true
				break
			}
		}
		if dominated {
			continue
		}
		// Candidate is not dominated, add to skyline and remove any skyline points it dominates
		var newSkyline []Point
		for _, s := range updatedSkyline {
			if utilities.Dominates(candidate, s, e.prefs) {
				continue
			}
			newSkyline = append(newSkyline, s)
		}
		newSkyline = append(newSkyline, candidate)
		updatedSkyline = newSkyline
	}
	e.skyline = updatedSkyline
}

// Skyline returns the current skyline set.
func (e *engine) Skyline() []Point {
	return e.skyline
}

// InsertBatch adds multiple new points and updates the skyline using the configured algorithm (default BNL).
// All new points are considered together with the current skyline, and only the non-dominated points are kept.
func (e *engine) InsertBatch(points []Point) {
	e.points = append(e.points, points...)
	candidates := append(append([]Point(nil), e.skyline...), points...)
	skyline, err := Skyline(candidates, e.dims, e.prefs, e.algo)
	if err != nil {
		// fallback: use BNL if the configured algorithm fails
		skyline, _ = Skyline(candidates, e.dims, e.prefs, "bnl")
	}
	e.skyline = skyline
}

// equalPoint compares two points for equality.
func equalPoint(a, b Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
