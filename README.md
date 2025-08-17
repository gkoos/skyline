# Skyline

[![Go Reference](https://pkg.go.dev/badge/github.com/gkoos/skyline.svg)](https://pkg.go.dev/github.com/gkoos/skyline)
[![Go Report Card](https://goreportcard.com/badge/github.com/gkoos/skyline)](https://goreportcard.com/report/github.com/gkoos/skyline)
[![codecov](https://codecov.io/gh/gkoos/skyline/branch/main/graph/badge.svg)](https://codecov.io/gh/gkoos/skyline)

A Go library for **skyline queries** — multi-dimensional optimization to find the set of Pareto-optimal points from a dataset.  

Supports both **static** and **dynamic** skyline computation with multiple algorithms.

---

## What is a Skyline Query?

Skyline queries identify the subset of points in a multi-dimensional dataset that are **not dominated** by any other point.  
A point *A* dominates *B* if *A* is as good or better than *B* in all dimensions and strictly better in at least one dimension.

This is useful in scenarios like:
- Finding products that are best across multiple criteria (price, performance, battery life)
- Multi-criteria decision making in finance, logistics, recommendation systems, etc.

The result, called the **skyline set**, represents the **Pareto-optimal front** of your data.

---

## What does this library do?

- Compute skyline points from static datasets using multiple algorithms (Block Nested Loop, Divide & Conquer, SkyTree)
- Support dynamic updates: insert, batch insertdelete, and update points incrementally without recomputing from scratch
- Allow flexible dimension selection and preference (minimize/maximize per dimension)
- Provide a simple, idiomatic Go API for both static and dynamic skyline queries

---
- **Block Nested Loop (BNL):** Brute-force, compares each point to all others. Use `BlockNestedLoop(points, epsilon)`. Epsilon controls the dominance threshold.
- **Divide and Conquer (DNC):** Recursive, efficient for large datasets. Use `DivideAndConquer(points, DNCConfig{Epsilon: ...})`. Epsilon is passed via the config.
- **SkyTree:** Tree-based, efficient for high dimensions. Use `SkyTree(points, epsilon)`. Epsilon is passed as a parameter.

## Installation

```bash
go get github.com/gkoos/skyline
```

---

## API Overview

### Types

```go
type Point map[string]float64
type Preference map[string]Order
type Order int

const (
    Min Order = iota
    Max
    Ignore  // Skip this dimension in dominance comparisons
)
```

### Static Computation

```go
func Skyline(points []Point, dims []string, prefs Preference, algo string) ([]Point, error)
```

Computes the skyline from a static dataset.
- `points`: input points
- `dims`: dimensions to consider
- `prefs`: preferences per dimension (Min or Max)
- `algo`: algorithm to use (`"bnl"`, `"dnc"`, `"skytree"`)

### Dynamic Updates

You can use the dynamic skyline engine for incremental and batch updates. Two constructors are available:

#### 1. DynamicSkyline (with initial skyline computation)

```go
engine, err := skyline.DynamicSkyline(points, dims, prefs, algo)
if err != nil {
    panic(err)
}
```
This computes the initial skyline from the dataset using the specified algorithm ("bnl", "dnc", or "skytree").

#### 2. DynamicSkylineRaw (no initial skyline computation)

```go
engine := skyline.DynamicSkylineRaw(points, dims, prefs, algo)
```
This skips initial skyline computation and treats the provided points as the current skyline. Useful for streaming or batch scenarios.

#### Supported Operations

- `engine.Insert(point)` — Insert a single point (uses BNL logic)
- `engine.InsertBatch(points)` — Insert multiple points at once (uses the configured algorithm for batch skyline computation)
- `engine.Update(oldPoint, newPoint)` — Replace a point and update the skyline
- `engine.Delete(point)` — Remove a point and update the skyline
- `engine.Skyline()` — Get the current skyline set

#### Example

```go
package main

import (
    "fmt"
    "github.com/yourname/skyline"
)

func main() {
    points := []skyline.Point{
        {"price": 400, "battery": 10},
        {"price": 500, "battery": 12},
        {"price": 300, "battery": 9},
        {"price": 450, "battery": 11},
    }

    prefs := skyline.Preference{
        "price": skyline.Min,
        "battery": skyline.Max,
    }

    // Dynamic skyline with initial computation
    engine, err := skyline.DynamicSkyline(points, []string{"price", "battery"}, prefs, "dnc")
    if err != nil {
        panic(err)
    }

    // Insert a single point
    newPoint := skyline.Point{"price": 420, "battery": 15}
    engine.Insert(newPoint)
    fmt.Println("After insert:", engine.Skyline())

    // Batch insert
    batch := []skyline.Point{
        {"price": 410, "battery": 13},
        {"price": 390, "battery": 16},
    }
    engine.InsertBatch(batch)
    fmt.Println("After batch insert:", engine.Skyline())

    // Update a point
    updatedPoint := skyline.Point{"price": 460, "battery": 14}
    engine.Update(newPoint, updatedPoint)
    fmt.Println("After update:", engine.Skyline())

    // Delete a point
    engine.Delete(updatedPoint)
    fmt.Println("After delete:", engine.Skyline())
}
```

### Partial Skyline

`Preference` includes an `Ignore` option to skip dimensions in dominance checks. This allows you to compute skylines based on a subset of dimensions, which can be useful in scenarios where some dimensions are not relevant.\
A practical example is adding a unique key to each point which then can be ignored in dominance checks.

## Algorithms

### Block Nested Loop (BNL)
- Simple, intuitive algorithm
- Compares each point with all others to find dominating points
- Works well for small datasets and supports incremental updates easily
- *In dynamic mode, we always use this algorithm* to insert a single point

### Divide & Conquer (D&C)
- Recursively divides data into smaller subsets, computes skylines, and merges results
- More efficient than BNL for larger datasets
- Static algorithm, dynamic extension is complex

### SkyTree
- Advanced algorithm using tree structures to prune comparisons
- Scales well with high-dimensional and large datasets
- Designed primarily for static datasets

#### Optimization Steps
The SkyTree implementation in this library includes several advanced optimizations for performance and scalability:
- **Advanced Pivot Selection:** Uses median or custom pivot selection to improve partitioning and pruning efficiency
- **Parallelization:** SkyTree uses parallelism in two main phases:
    - **Parallel Recursion:** When the number of partitions (regions) exceeds a threshold, recursive calls for each partition are executed in parallel using goroutines. This allows the algorithm to process different branches of the tree concurrently, greatly speeding up computation on multicore systems.
    - **Parallel Merge:** After recursion, the partial skylines from each partition are merged in parallel using a pairwise, multi-stage approach. At each stage, pairs of skylines are merged concurrently, reducing the total merge time to log₂(N) stages for N partitions.
    - **Worker Pool:** Both recursion and merge parallelism are managed by a configurable worker pool, which limits the number of concurrent goroutines to avoid oversubscription and maximize CPU efficiency. The pool size defaults to the number of available CPU cores, but can be tuned for your workload.
- **Dominance Caching:** Caches dominance checks between points to avoid redundant computations, reducing overall work
- **Slice Reuse:** Minimizes memory allocations by reusing slices in recursive calls and helpers
- **Custom Deduplication:** Uses a fast custom join for point keys, improving deduplication speed for large datasets
- **Configurable Recursion Depth:** Allows limiting recursion depth to prevent stack overflow and excessive computation; falls back to BNL if the limit is reached
- **Small Partition BNL Switch:** If the number of points in a partition is lsmall, SkyTree will use the Block Nested Loop (BNL) algorithm for that partition instead of recursing further. This optimization avoids SkyTree's overhead on small datasets, where BNL is typically faster, and can significantly improve performance for workloads with many small partitions. The threshold is tunable; see the configuration section for details.

These optimizations make SkyTree suitable for very large and high-dimensional datasets, balancing speed, memory usage, and accuracy.

### Approximate Skyline Queries

All skyline algorithms in this package support an **epsilon** parameter, which allows for approximate dominance. Epsilon is a non-negative float that relaxes the strictness of dominance comparisons:

- **Strict dominance** (`epsilon = 0`): A point `A` strictly dominates point `B` if all coordinates of `A` are less than or equal to those of `B`, and at least one is strictly less.
- **Epsilon dominance** (`epsilon > 0`): A point `A` epsilon-dominates point `B` if all coordinates of `A` are less than or equal to those of `B` plus `epsilon`, and at least one coordinate is strictly less by more than `epsilon`.

This is useful for handling floating-point imprecision or for applications where small differences are not significant.

Epsilon dominance checks can be combined with sampling and partitioning of the data to calculate approximate skylines which can be more efficient for large datasets - essentially a tradeoff between accuracy and performance.

#### Example: Dominance with Epsilon

```go
import "skyline/internal/algorithms"

a := []float64{1.0, 2.0}
b := []float64{1.0, 2.0000001}

// Strict dominance (epsilon = 0): false
algorithms.DominatesEpsilon(a, b, 0) // false

// Epsilon dominance (epsilon = 1e-6): true
algorithms.DominatesEpsilon(a, b, 1e-6) // true
```

---

## Configuration

Skyline algorithms can be fine-tuned using configuration options to optimize performance and scalability for different dataset sizes and characteristics. Tuning these options allows you to balance speed, memory usage, and accuracy, especially for large or high-dimensional data. Proper configuration is essential to avoid bottlenecks, excessive memory consumption, or incomplete results, and lets you adapt the algorithms to your specific workload and hardware.

### Block Nested Loop (BNL)
- `Epsilon`: Dominance threshold for comparisons. Lower values increase accuracy but may slow down performance. Higher values speed up comparisons but may miss some dominated points. Default is `0.0`, meaning exact dominance checks.

### Divide & Conquer (DNC)
- `Threshold`: Minimum number of points in a partition before switching to BNL. Lower values increase recursion, higher values use BNL more often. Tune for your dataset size.
- `BatchSize`: Number of points processed together in each batch. Larger batches can improve cache locality and throughput, but may use more memory.
- `Epsilon`: Dominance threshold for comparisons. See Block Nested Loop (BNL) section.

### SkyTree
- `PivotSelector`: Function to choose the pivot point for partitioning. The default is median selection, but you can provide a custom function for domain-specific optimization.
- `ParallelThreshold`: Minimum number of partitions before enabling parallel processing. Lower values increase parallelism, higher values reduce goroutine overhead.
- `MaxRecursionDepth`: Maximum allowed recursion depth. If exceeded, SkyTree falls back to BNL for the remaining data. Prevents stack overflow and excessive computation for very large or complex datasets.
- `BNLSwitchThreshold`: If the number of points in a partition is less than or equal to this threshold, SkyTree will use the Block Nested Loop (BNL) algorithm for that partition instead of recursing further. This improves performance by avoiding SkyTree's overhead on small datasets, where BNL is typically faster. The default is 32, but you can tune this value for your workload and hardware. Lower values reduce BNL usage; higher values make SkyTree switch to BNL more often for small partitions.
- `WorkerPoolSize`: Controls the maximum number of goroutines (workers) used for parallel recursion and merging in SkyTree. Setting this to `0` (the default) will use the number of available CPU cores on your system, which is usually optimal for most workloads. You can set a specific positive value to limit CPU usage or experiment with different levels of parallelism. Increasing this value may improve performance on large, partitionable datasets, but setting it too high can cause oversubscription and reduce efficiency. For most users, leaving it at `0` is recommended.
- `Epsilon`: Dominance threshold for comparisons. See Block Nested Loop (BNL) section.

Refer to the code and examples for how to set these options in your application.

---

## Running Tests

To run all unit tests for algorithms and utilities:

```bash
go test ./...
```

This will execute all correctness and edge case tests for the skyline algorithms, domination logic, and dynamic updates. Tests cover small, large, and high-dimensional datasets, as well as pathological cases.

---

## Benchmarking Skyline Algorithms

To run performance benchmarks and compare algorithms:

```bash
go test -bench . ./internal/algorithms
```

This command runs repeatable benchmarks for all implemented skyline algorithms (BNL, D&C, SkyTree) on a variety of datasets, including small, large, high-dimensional, clustered, and pathological cases.

**Algorithm selection guidance:**
- **Block Nested Loop (BNL):** Best for small datasets or when incremental updates are needed. Simple, but slow for large or high-dimensional data.
- **Divide & Conquer (DNC):** Generally fastest for large, diverse datasets with a moderate skyline size. Uses more memory, but scales well unless most points are dominated.
- **SkyTree:** Optimized for very large, high-dimensional datasets with many dominated (clustered) points and a small skyline. SkyTree is much faster and more memory-efficient than DNC when the dataset is highly clustered and the skyline is small. For datasets where most points are Pareto-optimal (large skyline), DNC may be faster.

**Summary:**
- Use BNL for small or dynamic datasets.
- Use DNC for large, diverse datasets with a moderate skyline.
- Use SkyTree for large, high-dimensional, clustered datasets with a small skyline.
- Benchmark your own data to select the best algorithm for your use case.

---

## Future Improvements

- Add debug visualizer (especially for the SkyTree algorithm) (optional CLI)
- Add batch processing support to the SkyTree implementation

---

## License

MIT License

---

## Contact

Contributions and issues are welcome! Please open an issue or submit a pull request.

For questions, issues, or contributions, please use the GitHub repository:

https://github.com/gkoos/skyline
