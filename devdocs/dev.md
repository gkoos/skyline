# Folder Structure
```
skyline/
├── go.mod
├── go.sum
├── README.md
├── cmd/
│   └── demo/                   # Optional CLI or playground demo
│       └── main.go
├── internal/
│   └── bnl/                    # Block Nested Loop algorithm
│       └── bnl.go
│   └── dnc/                    # Divide and Conquer algorithm
│       └── dnc.go
│   └── skytree/                # SkyTree algorithm
│       └── skytree.go
│   └── common/                 # Shared utilities across algorithms
│       ├── dominance.go
│       ├── point.go
│       └── helpers.go
├── skyline/
│   ├── skyline.go              # Public API (static Skyline function)
│   ├── engine.go               # Dynamic engine interface and implementation
│   ├── types.go                # Point, Preference, Order, etc.
│   └── registry.go             # Algorithm registry / dispatcher
├── tests/
│   ├── skyline_test.go
│   ├── engine_test.go
│   └── dataset_generator.go   # Helper to generate synthetic test data
└── examples/
    └── basic_usage/
        └── main.go
```

`cmd/demo/`
A CLI or small playground app to test the API or debug your algorithms. Can evolve into a benchmark suite or devtool.

`internal/`
Houses internal algorithm implementations (BNL, D&C, SkyTree), kept separate from public APIs. Also includes shared logic like point comparison, dominance checks, etc.

`skyline/`
Public package that exposes a clean interface:
- Skyline() function for static computation.
- Engine interface for dynamic use cases.
- Common types like Point, Preference, etc.
This package depends on internal/ but is what users will import.

`tests/`
Holds both unit and integration tests. You can include test data generators for benchmarking or fuzzing.

`examples/`
Contains real-world usage examples with different types of data and configurations, useful for docs and community adoption.

