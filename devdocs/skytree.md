# Skytree algorithm

## What SkyTree Actually Does

The SkyTree algorithm:
- Recursively splits the dataset into subspaces using dominance regions.
- In each subspace, it prunes dominated points early.
- Maintains a tree structure to represent dominance relationships, hence the name SkyTree.
- Uses bitstrings or binary vectors to track dominance regions efficiently.

## Key Concepts

### 1. Dominance

Point A dominates point B if:
- A is no worse in all dimensions, and
- A is better in at least one dimension.

### 2. Region Bitstring

Every point can be associated with a region code based on how it compares to a reference point in each dimension:
- 0 if it's equal
- 1 if it's better (lower)
- 2 if it's worse (higher)

Example:
If point p1 = [3, 5] and p2 = [4, 4], then comparing p2 to p1:

    In dimension 1: 4 > 3 → 2

    In dimension 2: 4 < 5 → 1
    So region bitstring = [2, 1]

This region is used to group points in recursive partitions.


###3. Pruning

Points in certain regions can’t dominate others — they can be skipped early. For example, if a region is fully dominated by an already known skyline point, skip it.

## Step-by-Step SkyTree Algorithm

Recursive Function

```go
SkyTree(points, dimensions):
    if points is empty or has only 1 point:
        return points

    choose a pivot point (usually median or random)
    create a hash map of region → points[]
    
    for each point in points (excluding pivot):
        compute region_bitstring = compare(point, pivot)
        assign to region[region_bitstring]

    S = {pivot}

    for region in dominance-prunable order:
        if region can be pruned based on dominance:
            skip
        else:
            subtree_skyline = SkyTree(region_points, dimensions)
            prune subtree_skyline using S
            add non-dominated points from subtree_skyline to S

    return S
```

## Practical Implementation Strategy

### 1. Preprocessing

- Normalize data (optional, for real-world cases)
- Sort or cluster points to improve pruning efficiency

### 2. Region Encoding

Implement a method to generate region bitstrings for each point compared to a pivot.

### 3. Prune Logic

You’ll need a helper to:
- Check if a region can be pruned
- Remove dominated points from partial skylines

### 4. SkyTree Building

Implement the recursive structure:
- Base case
- Region splitting
- Recursive call per region

## Suggested Helper Functions

| Function	| Purpose |
| -----     | ----- |
| dominates(a, b)	| Returns true if a dominates b
| compareRegion(p, pivot)	| Returns bitstring of the region
| prune(region, knownSkyline) |	Filters out dominated points
| canPruneRegion(regionBitstring, pivotRegion)	| Logical condition for skipping a region

## Optimization Notes

- Use bit masking for region encoding
- Consider parallelizing the recursive region calls (ideal for Go)
- Cache dominance results when possible
- Limit recursion depth with early pruning



## ChatGPT implementation
```go
// Package skytree implements the SkyTree algorithm for skyline queries
// with full support for bitmask lattice partitioning, advanced pivot selection,
// and high-dimensional pruning.
package skytree

import (
	"math/rand"
	"sort"
)

type Point struct {
	Coords []int
}

func (p Point) Dominates(q Point) bool {
	better := false
	for i := range p.Coords {
		if p.Coords[i] > q.Coords[i] {
			return false
		}
		if p.Coords[i] < q.Coords[i] {
			better = true
		}
	}
	return better
}

func (p Point) DominationMask(pivot Point) int {
	mask := 0
	for i := range p.Coords {
		if p.Coords[i] < pivot.Coords[i] {
			mask |= 1 << i
		}
	}
	return mask
}

// Heuristic pivot selection: choose the point with smallest sum of coordinates
func SelectPivot(points []Point) Point {
	minSum := int(^uint(0) >> 1) // Max int
	best := points[0]
	for _, p := range points {
		s := 0
		for _, v := range p.Coords {
			s += v
		}
		if s < minSum {
			minSum = s
			best = p
		}
	}
	return best
}

type SkyTreeNode struct {
	Pivot    Point
	Region   int
	Children []*SkyTreeNode
}

func BuildSkyTree(points []Point, dim int) *SkyTreeNode {
	if len(points) == 0 {
		return nil
	}
	if len(points) == 1 {
		return &SkyTreeNode{Pivot: points[0]}
	}

	pivot := SelectPivot(points)
	remaining := []Point{}

	// Remove points dominated by pivot, and improve pivot if needed
	for _, pt := range points {
		if pivot.Dominates(pt) {
			continue
		}
		if pt.Dominates(pivot) {
			pivot = pt
			continue
		}
		remaining = append(remaining, pt)
	}

	// Partition points by lattice region
	partitions := make(map[int][]Point)
	for _, pt := range remaining {
		mask := pt.DominationMask(pivot)
		partitions[mask] = append(partitions[mask], pt)
	}

	node := &SkyTreeNode{
		Pivot:    pivot,
		Children: []*SkyTreeNode{},
	}

	// Order partitions by mask complexity (popcount)
	sortedMasks := make([]int, 0, len(partitions))
	for m := range partitions {
		sortedMasks = append(sortedMasks, m)
	}
	sort.Slice(sortedMasks, func(i, j int) bool {
		return popCount(sortedMasks[i]) < popCount(sortedMasks[j])
	})

	// Build subtrees with pruning
	for _, mask := range sortedMasks {
		subset := partitions[mask]
		if len(subset) == 0 {
			continue
		}
		pruned := pruneWithPivot(subset, pivot)
		child := BuildSkyTree(pruned, dim)
		if child != nil {
			child.Region = mask
			node.Children = append(node.Children, child)
		}
	}

	return node
}

func pruneWithPivot(points []Point, pivot Point) []Point {
	res := []Point{}
	for _, p := range points {
		if !pivot.Dominates(p) {
			res = append(res, p)
		}
	}
	return res
}

func popCount(n int) int {
	count := 0
	for n > 0 {
		count += n & 1
		n >>= 1
	}
	return count
}

func CollectSkyline(node *SkyTreeNode, result *[]Point) {
	if node == nil {
		return
	}
	*result = append(*result, node.Pivot)
	for _, child := range node.Children {
		CollectSkyline(child, result)
	}
}

func Skyline(points []Point, dim int) []Point {
	tree := BuildSkyTree(points, dim)
	result := []Point{}
	CollectSkyline(tree, &result)
	return result
}
```

