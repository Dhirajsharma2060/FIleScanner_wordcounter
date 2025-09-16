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
+-------------------+
|   Worker Pool     |  (Consumers)
+-------------------+
          |
          v
+-------------------+
|   Main Goroutine  |  (Aggregation & Top-K)
+-------------------+
          |
          v
+-------------------+
|     Output        |
+-------------------+
```

---

## 🔄 Concurrency Model

- **Producer-Consumer Pattern:**  
  - The directory walker (producer) streams file paths to a buffered channel.
  - Multiple worker goroutines (consumers) process files concurrently.

- **Direct Aggregation:**  
  - Each worker sends its local word frequency map to a single results channel (`resultChan`).
  - The main goroutine merges all results into a global map.

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
  - Sends the entire local frequency map to `resultChan`

- **Main Goroutine:**  
  - Receives all local frequency maps from workers via `resultChan`
  - Merges them into a single `wordFreq` map
  - Calculates and prints the top-K words

---

## 🛡️ Error Handling

- All file and directory errors are reported and do not silently fail.
- The directory walk returns an error if traversal fails.
- Worker errors are printed with context.
- Context cancellation is handled for graceful shutdown on SIGINT/SIGTERM.

---

## 🏎️ Performance & Scalability

- **Parallelism:**  
  - Number of workers = number of CPU cores (configurable)
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

- Easily configurable for directory, top-K, and worker count.
- Can be extended to support other file types or more advanced text processing.

---

## 📚 References

- See [README.md](README.md) for usage and features.
- See [specification.txt](specification.txt) for assignment