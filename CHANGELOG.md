# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2025-08-18

### Added
- Epsilon parameter for approximate dominance in all skyline algorithms (static and dynamic)
- Batch insert support in dynamic skyline engine

### Changed
- README updated for epsilon and batch features

### Fixed
- Improved test coverage for dynamic and static skyline operations

## [1.2.0] - 2025-08-07

### Added
- Partial skyline computation with `Ignore` preference type
- Configured golangci-lint for Go 1.24 compatibility with custom linter selection

### Changed
- Refactored `skytreeRecWithDepthPool` function to reduce cyclomatic complexity from 24 to 5
- Updated golangci-lint configuration to work around Go 1.24 compatibility issues

### Fixed
- Applied `gofmt -s` formatting across entire codebase
