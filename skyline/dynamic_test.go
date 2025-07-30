package skyline

import (
	"testing"
)

// 5000 points, a couple dominating
func makeDataset5000CoupleDominating() Dataset {
	data := make(Dataset, 5000)
	// 4997 points clustered in the lower left
	for i := 0; i < 4997; i++ {
		data[i] = Point{float64(i % 50), float64(i % 50)}
	}
	// 3 dominating points in the upper right
	data[4997] = Point{1000, 1000}
	data[4998] = Point{900, 1100}
	data[4999] = Point{1100, 900}
	return data
}

func TestDynamicInsertNonSkyline(t *testing.T) {
	engine, err := DynamicSkyline(makeDataset5000CoupleDominating(), []string{"0", "1"}, Preference{Max, Max}, "bnl")
	if err != nil {
		t.Fatalf("engine creation failed: %v", err)
	}
	before := engine.Skyline()
	// Add a point that is dominated by the skyline
	engine.Insert(Point{500, 500})
	after := engine.Skyline()
	if len(before) != len(after) {
		t.Errorf("Skyline changed after inserting non-skyline point")
	}
}

func TestDynamicInsertSkyline(t *testing.T) {
	engine, err := DynamicSkyline(makeDataset5000CoupleDominating(), []string{"0", "1"}, Preference{Max, Max}, "bnl")
	if err != nil {
		t.Fatalf("engine creation failed: %v", err)
	}
	before := engine.Skyline()
	// Add a point that should be part of the skyline
	engine.Insert(Point{-10, 2000})
	after := engine.Skyline()
	if len(after) == len(before) {
		t.Errorf("Skyline did not change after inserting skyline point")
	}
}

func TestDynamicInsertDominatingAll(t *testing.T) {
	engine, err := DynamicSkyline(makeDataset5000CoupleDominating(), []string{"0", "1"}, Preference{Max, Max}, "bnl")
	if err != nil {
		t.Fatalf("engine creation failed: %v", err)
	}
	// Add a point that dominates all others
	engine.Insert(Point{2000, 2000})
	after := engine.Skyline()
	if len(after) != 1 || !equalPoint(after[0], Point{2000, 2000}) {
		t.Errorf("Skyline not replaced by dominating point")
	}
}

func TestDynamicDeleteNonSkyline(t *testing.T) {
	engine, err := DynamicSkyline(makeDataset5000CoupleDominating(), []string{"0", "1"}, Preference{Max, Max}, "bnl")
	if err != nil {
		t.Fatalf("engine creation failed: %v", err)
	}
	before := engine.Skyline()
	engine.Delete(Point{5000, 500})
	after := engine.Skyline()
	if len(before) != len(after) {
		t.Errorf("Skyline changed after deleting non-skyline point")
	}
}

func TestDynamicDeleteSkyline(t *testing.T) {
	engine, err := DynamicSkyline(makeDataset5000CoupleDominating(), []string{"0", "1"}, Preference{Max, Max}, "bnl")
	if err != nil {
		t.Fatalf("engine creation failed: %v", err)
	}
	before := engine.Skyline()
	engine.Delete(Point{1000, 1000})
	after := engine.Skyline()
	if len(after) == len(before) {
		t.Errorf("Skyline did not change after deleting skyline point")
	}
}
