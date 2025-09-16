# System Design & Architecture: FileScanner WordCounter

This document describes the system architecture, concurrency model, and design decisions behind the FileScanner WordCounter project.

---

## 🏗️ System Overview

The tool is a concurrent, pipeline-based CLI application written in Go. It efficiently scans directories for `.txt` files, counts lines, and computes the top-K most frequent words across all files.

---

## ⚙️ Architecture Diagram

```
+-------------------+
| Directory Walker  |  (Producer)
+-------------------+
          |
          v
+-------------------+
|   fileChan (chan) |
+-------------------+
          |
          v
+-------------------+      +-------------------+
|   Worker Pool     | ---> |   Aggregator Pool |
| (Consumers)       |      | (Sharded by hash) |
+-------------------+      +-------------------+
          |                        |
          +------------------------+
                   |
                   v
           +-----------------+
           |   Final Merge   |
           +-----------------+
                   |
                   v
           +-----------------+
           |   Top-K Heap    |
           +-----------------+
                   |
                   v
           +-----------------+
           |     Output      |
           +-----------------+
```

---

## 🔄 Concurrency Model

- **Producer-Consumer Pattern:**  
  - The directory walker (producer) streams file paths to a buffered channel.
  - Multiple worker goroutines (consumers) process files concurrently.

- **Sharded Aggregation:**  
  - Each worker sends its local word frequency map to one of several aggregator goroutines, chosen by hashing the file name.
  - Each aggregator merges results into its own map, reducing contention.

- **Final Merge:**  
  - After all aggregators finish, their maps are merged into a single global map.

- **Top-K Calculation:**  
  - A min-heap is used to efficiently extract the top K most frequent words.

---

## 🧩 Key Components

- **Directory Walker:**  
  Recursively scans for `.txt` files and sends their paths to `fileChan`.

- **Worker Pool:**  
  Each worker:
  - Reads file paths from `fileChan`
  - Counts lines and word frequencies (with normalization)
  - Sends results to the appropriate aggregator

- **Aggregator Pool:**  
  Each aggregator:
  - Receives word counts from workers (for its shard)
  - Merges them into a local map

- **Final Aggregation:**  
  Merges all aggregator maps into a single `wordFreq` map.

- **Top-K Heap:**  
  Maintains the top K most frequent words using a min-heap for efficiency.

---

## 🛡️ Error Handling

- All file and directory errors are reported and do not silently fail.
- The directory walk returns an error if traversal fails.
- Worker errors are printed with context.

---

## 🏎️ Performance & Scalability

- **Parallelism:**  
  - Number of workers = number of CPU cores (configurable)
  - Number of aggregators = 4 (configurable)
- **Memory Efficiency:**  
  - Files are streamed, not preloaded.
  - Only top K words are kept in the heap.
- **No Global Mutable State:**  
  - All communication is via channels and local maps.

---

## 🧪 Testing

- Unit tests cover normalization, file processing, heap logic, and directory walking.

---

## 🔧 Extensibility

- Easily configurable for directory, top-K, worker count, and aggregator count.
- Can be extended to support other file types or more advanced text processing.

---

## 📚 References

- See [README.md](README.md) for usage and features.
- See [specification.txt](specification.txt) for assignment