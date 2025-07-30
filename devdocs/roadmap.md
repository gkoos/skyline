# 🛣️ Skyline Query Engine — Project Roadmap

This project implements efficient static and dynamic skyline query algorithms in Go, including the SkyTree algorithm.

---

## ✅ Phase 0: Foundation & Planning

| Task | Description |
|------|-------------|
| 🔧 Define Scope | Implement static skyline queries (BNL, D&C, SkyTree) + basic dynamic support (insert/update/delete). |
| 📁 Design Folder Structure | Set up `pkg/`, `internal/`, `cmd/`, and test directories. |
| 🧪 Define Use Cases | E.g., product filtering, multi-criteria optimization, real-time dynamic datasets. |
| 📚 Draft README | Include what skyline queries are, intended audience, and API vision. |

---

## 🧱 Phase 1: Static Skyline Algorithms

### 1.1: Baseline Algorithm – BNL (Block Nested Loops)
- [x] `Dominates(a, b)` utility function
- [x] Implement naive BNL algorithm
- [x] Add unit tests + benchmarks

### 1.2: Divide and Conquer (D&C)
- [x] Recursive splitting and merging
- [x] Efficient dominance merge logic
- [x] Benchmark vs BNL

---

## 🔁 Phase 2: Dynamic Skyline Support

### 2.1: Incremental Updates
- [x] Skyline insert: check if new point dominates or is dominated
- [x] Skyline delete: reprocess shadowed region
- [x] Skyline update: treat as delete + insert


---

## 🧪 Phase 3: API & Testing

| Task | Description |
|------|-------------|
| 🔌 Public API | Finalize public methods: `skyline.Static(points)`, `skyline.Dynamic().Insert(...)`, etc. |
| 🧪 Unit Tests | Cover edge cases, corner cases, and performance tests |
| 📈 Benchmarking Suite | Use Go’s benchmarking tools across algorithms |
| ✅ Fuzz Testing | Use `testing/quick` for randomized input tests |

---

## 📚 Phase 4: SkyTree Algorithm Implementation

### 4.1: SkyTree
- [ ] Region encoding & dominance rules
- [ ] Region-to-subspace grouping
- [ ] Recursively build sky tree
- [ ] Add debug visualizer (optional CLI)
- [ ] Benchmark against BNL and D&C

---

## 🌐 Phase 5: Documentation & Examples

| Task | Description |
|------|-------------|
| 📖 Full API Docs | Add GoDoc documentation |
| 📦 Examples | Product filter, multi-criteria selection, live update handling |
| 📘 Blog / Medium Post | Walkthrough of skyline queries and your implementation |

---

## 🚀 Phase 6: Release & Community

| Task | Description |
|------|-------------|
| 🔖 v0.1.0 Release | Publish on GitHub and pkg.go.dev |
| 📣 Launch Post | Reddit / Hacker News / Golang Subreddit |
| 🔄 Gather Feedback | Track bugs, performance suggestions |
| 🧩 Explore Bindings | Optional: WASM + Node bindings if real-time JS clients are needed |
