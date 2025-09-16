# FileScanner WordCounter

A high-performance, concurrent command-line tool in Go for scanning directories, counting lines in `.txt` files, and reporting the top K most frequent words across all files.

---

## ✨ Features

- **Recursive Directory Scan:** Finds all `.txt` files in a directory and its subdirectories.
- **Parallel Processing:** Uses multiple goroutines to process files concurrently for maximum speed.
- **Line Counting:** Prints the filename and the total number of lines for each `.txt` file.
- **Word Frequency Analysis:** Counts word frequencies across all files, normalizing words (case-insensitive, punctuation stripped).
- **Top-K Words:** Efficiently reports the top 10 most frequent words using a min-heap (O(N log K) time).
- **Robust Error Handling:** Handles file and directory errors gracefully.
- **Scalable Aggregation:** Uses multiple aggregator goroutines to avoid bottlenecks during result merging.
- **Unit Tested:** Includes unit tests for core logic.

---

## 🚀 Usage

### 1. Build & Run

```sh
go run [main.go](http://_vscodecontentref_/0) <directory>
```
### 2. Output

- **For each .txt file:** 
```sh
[Worker X] File: <filename>, Lines: <count>
```
- **After processing:**
- ***Aggregation time***
- ***Top 10 words and their frequencies***
- ***Top words calculation time***
- ***Total execution time***




## 🛠️ How It Works

1. **Directory Traversal:**
Recursively walks the given directory, sending .txt file paths to a channel.

2. **Producer-Consumer Pattern:** 

- ***Producer:*** Directory walker sends file paths to a channel.
- ***Consumers:*** Worker goroutines read file paths, process files, and count word frequencies.

3. **Concurrent Aggregation:**

- Each worker sends its local word frequency map to one of several aggregator goroutines (sharded by hash).
- Aggregators merge results in parallel, then the main goroutine combines them for the final tally.
4. **Top-K Calculation:**
- Uses a min-heap to efficiently keep only the top 10 most frequent words.
5. **Word Normalization:**
- Words are lowercased and stripped of punctuation for accurate counting.

## 🧪 Testing

```
go test
```
## 📁 Project Structure

```
.
├── [go.mod](http://_vscodecontentref_/2)
├── [main.go](http://_vscodecontentref_/3)
├── [main_test.go](http://_vscodecontentref_/4)
├── [README.md](http://_vscodecontentref_/5)
├── [specification.txt](http://_vscodecontentref_/6)
├── [test1.txt](http://_vscodecontentref_/7)
├── [test2.txt](http://_vscodecontentref_/8)

```
## 📋 Example Output
``` 
[Worker 2] File: test2.txt, Lines: 3
[Worker 0] File: test1.txt, Lines: 584
Found 2 text files. Using 8 workers.
Aggregation time: 1.065311ms
Top 10 words:
the: 1033
of: 465
and: 357
to: 240
in: 212
a: 185
is: 133
by: 104
on: 103
mahabharata: 90
Top words calculation time: 648.567µs
Total execution time: 6.763312ms

```

## 📖 Specification
See specification.txt for the full assignment requirements.

## 🙏 Acknowledgements

- Inspired Wikipedia text samples.
- Built with Go’s powerful concurrency primitives.