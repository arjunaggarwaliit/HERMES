# COSMOS: Concurrent Optimal Shortest Path and Modular Update System

COSMOS is a system for computing shortest paths and processing real time updates on large, dynamic road networks. It combines hierarchical graph partitioning, a shortcut network overlay, and a concurrency control scheme so that routing queries and network updates can run at the same time without blocking each other or producing inconsistent results.

This repository contains multiple implementations of the core algorithms (Go, Java, C++), a vendored Rust baseline used for comparison, datasets used for testing, and the analysis scripts used to produce the performance results described in the accompanying report.

## Background

Classic shortest path computation on road networks relies on Dijkstra's algorithm or its bidirectional variant. These are correct but slow on continental scale graphs, which is why most practical routing systems preprocess the graph to add shortcuts that let queries skip over large sections of the network. Well known approaches in this space include ALT, Contraction Hierarchies, and Arc Flags.

Static preprocessing does not handle networks that change over time, for example due to live traffic conditions or temporary closures. COSMOS is built around the idea that a routing system needs to support both fast queries and frequent updates at the same time, which means the query and update paths need explicit concurrency control rather than assuming the graph is fixed.

## How COSMOS works

### 1. Partitioning: Natural Cuts Based Connectivity Clustering (NCCC)

The graph is broken into partitions in three phases:

- **Coarsening**: the graph is progressively simplified to reduce its size while preserving connectivity structure, guided by a weighted connectivity/reduction ratio.
- **Cutting**: partitions are split using a two stage process, first identifying candidate single cuts and then refining them into natural, low cost cuts along boundaries that minimize the number of edges crossing partitions.
- **Integration**: the refined partitions are merged back into a consistent structure, producing fragments that form the base level of the hierarchy.

### 2. Multi-Level Partitioning (MLP) and the shortcut network

Partitions are organized into a multi-level hierarchy. Within each partition, nodes are classified based on their role (interior vs. border nodes), and an all-pairs shortest path computation is run to build a shortcut network connecting border nodes. These shortcuts are overlaid onto the level above, so higher levels of the hierarchy only need to reason about a much smaller graph made of shortcuts rather than the full road network.

### 3. Routing queries: Hierarchical Bidirectional Dijkstra

A routing query runs a forward search from the source and a backward search from the destination simultaneously. Both searches climb the partition hierarchy toward border nodes rather than exploring the full graph. When the two searches meet at the lowest common ancestor partition, the path length is known immediately. The full path is then recovered with a path unpacking step, which recursively expands each shortcut back into the underlying road segments.

### 4. Update processing

Updates to the network (new edges, changed weights, closures) are batched over a fixed interval called the Delta Phase rather than applied one at a time. At the end of each interval, affected partitions are identified, their internal shortest paths are recomputed, and the partition hierarchy is reused wherever connectivity has not changed. This keeps update cost proportional to the size of the affected region rather than the whole graph.

### 5. Concurrency control

COSMOS uses the Delta Phase batching described above together with Two-Phase Locking (2PL) to guarantee serializable execution when queries and updates run concurrently. Locks are acquired during a growing phase and released during a shrinking phase, which is the standard mechanism for guaranteeing conflict serializability in database systems, applied here to graph partitions instead of database rows.

##  Repository Structure

```text
COSMOS/
├── go-spcs/                         Main Go implementation of COSMOS
│   ├── spcs/                        Entry point and example programs
│   ├── dch/                         Core graph processing and routing engine
│   │   ├── utils/                   Graph structures, Dijkstra, contraction,
│   │   │                            import/export utilities, heaps, etc.
│   │   └── tests/                   Unit tests for core algorithms
│   └── demo/                        Experimental framework and evaluation suite
│       ├── src/                     Partitioning, MLP, routing, update processing,
│       │                            and concurrency control
│       ├── simulator/               Workload and concurrency simulators
│       ├── runnable/                Executable entry points and benchmark configurations
│       ├── analyzers/               Python scripts for log analysis and visualization
│       └── tests/                   End-to-end integration tests
│
├── java-spcs/                       Java implementation (Gradle project)
│   └── app/src/main/java/csps/      Graph data structures and algorithm implementations
│
├── cpp-spcs/                        Early C++ prototype and experimental implementation
│   ├── src/                         Source code and headers
│   ├── test/                        Sample graph datasets
│   └── plot/                        Gnuplot scripts and benchmark visualizations
│
├── rust-spcs/                       Vendored implementation of the "fast_paths" library
│                                    (MIT/Apache-2.0) used as the Contraction
│                                    Hierarchies baseline for benchmarking
│
├── datasets/                        Road network datasets and preprocessing utilities
│
└── go.work                          Go workspace configuration for the Go modules
```
## Getting started

### Go implementation (primary)

Requires Go 1.21.6 or later.

```bash
cd COSMOS
go work sync
cd go-spcs/demo/runnable
go run main.go
```

Run configuration, including which graph file to load, the number of hierarchy levels, and which test suite to run (routing only, or routing plus concurrent updates), is set in `go-spcs/demo/config/config.yaml`.

To run the Go unit and integration tests:

```bash
cd go-spcs/dch/tests
go test ./...

cd ../../demo/tests
go test ./...
```

### Java implementation

Requires a JDK compatible with the Gradle wrapper included in the repository.

```bash
cd java-spcs
./gradlew build
./gradlew test
```

### C++ implementation

```bash
cd cpp-spcs/src
make
./edp_main
```

### Rust baseline (Contraction Hierarchies comparison)

```bash
cd rust-spcs
cargo build --release
```

### Datasets

`datasets/dataset_converter.py` converts raw graph data into the CSV/edge list format expected by the Go, Java, and C++ implementations. Sample datasets are provided under `datasets/misc/`.

## Performance summary

The system was evaluated against Contraction Hierarchies (CH) as a baseline on the same road network.

**Correctness**: verified with a confusion matrix comparing expected versus computed query costs across a large batch of routing queries, with results aligning along the diagonal as expected for a correct implementation.

**Routing query latency (microseconds), by percentile**

| Percentile | CH   | COSMOS |
|------------|------|--------|
| 50th       | 0.0  | 0.0    |
| 90th       | 1.00 | 0.38   |
| 99th       | 6.99 | 1.15   |

**Update query latency (milliseconds), by percentile**

| Percentile | COSMOS |
|------------|--------|
| 50th       | 38.45  |
| 99th       | 67.50  |

COSMOS showed lower tail latency on routing queries than the CH baseline, and stable update times even under concurrent routing load, since updates are batched and query execution is not blocked while a batch is being applied.

## Third party components

`rust-spcs` bundles the `fast_paths` crate by easbar (https://github.com/easbar/fast_paths), licensed under MIT/Apache-2.0, and is used only as the Contraction Hierarchies baseline for benchmarking. It is not part of the COSMOS implementation itself; see `rust-spcs/LICENSE-MIT` and `rust-spcs/LICENSE-APACHE` for its license terms.

