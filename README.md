# Skyline

A Go library for **skyline queries** — multi-dimensional optimization to find the set of Pareto-optimal points from a dataset.  
Supports both **static** and **dynamic** skyline computation with multiple algorithms.

---

## What is a Skyline Query?

Skyline queries identify the subset of points in a multi-dimensional dataset that are **not dominated** by any other point.  
A point *A* dominates *B* if *A* is as good or better than *B* in all dimensions and strictly better in at least one dimension.  

This is useful in scenarios like:

- Finding products that are best across multiple criteria (price, performance, battery life).  
- Multi-criteria decision making in finance, logistics, recommendation systems, etc.

The result, called the **skyline set**, represents the **Pareto-optimal front** of your data.

---

## What does this library do?

- Compute skyline points from static datasets using multiple algorithms (Block Nested Loop, Divide & Conquer, SkyTree).  
- Support dynamic updates: insert, delete, and update points incrementally without recomputing from scratch.  
- Allow flexible dimension selection and preference (minimize/maximize per dimension).  
- Provide a simple, idiomatic Go API for both static and dynamic skyline queries.

---

## Installation

```bash
go get github.com/gkoos/skyline
```

## API Overview

### Types

```go
type Point map[string]float64
type Preference map[string]Order
type Order int

const (
    Min Order = iota
    Max
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

    // Static skyline
    result, err := skyline.Skyline(points, []string{"price", "battery"}, prefs, "dnc")
    if err != nil {
        panic(err)
    }
    fmt.Println("Static skyline:", result)

    // Dynamic skyline
    engine, err := skyline.DynamicSkyline(points, []string{"price", "battery"}, prefs, "dnc")
    if err != nil {
        panic(err)
    }
    newPoint := skyline.Point{"price": 420, "battery": 15}
    engine.Insert(newPoint)
    fmt.Println("After insert:", engine.Skyline())

    updatedPoint := skyline.Point{"price": 460, "battery": 14}
    engine.Update(newPoint, updatedPoint)
    fmt.Println("After update:", engine.Skyline())

    engine.Delete(updatedPoint)
    fmt.Println("After delete:", engine.Skyline())
}
```

## Algorithms

### Block Nested Loop (BNL)

- Simple, intuitive algorithm.
- Compares each point with all others to find dominating points.
- Works well for small datasets and supports incremental updates easily.
- *In dynamic mode, we always use this algorithm.*

### Divide & Conquer (D&C)

- Recursively divides data into smaller subsets, computes skylines, and merges results.
- More efficient than BNL for larger datasets.
- Static algorithm, dynamic extension is complex.

### SkyTree

- Advanced algorithm using tree structures to prune comparisons.
- Scales well with high-dimensional and large datasets.
- Designed primarily for static datasets.

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

This runs repeatable benchmarks for BNL and D&C on various datasets (small, large, high-dimensional, and pathological).
D&C is generally faster for large, diverse datasets; BNL may be faster for small or highly uniform data.

Use these results to select the best algorithm for your use case and to validate performance improvements.

## License

MIT License

## Contact

Contributions and issues are welcome! Please open an issue or submit a pull request.

For questions, issues, or contributions, please use the GitHub repository:

https://github.com/gkoos/skyline
