package main

import (
	"bufio"
	"container/heap"
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
)

func main() {
	walkStart := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory path")
		return
	}
	root := os.Args[1]

	numWorkers := runtime.NumCPU()
	fileChan := make(chan string, 100)
	var wg sync.WaitGroup

	numAggregators := 4
	aggChans := make([]chan map[string]int, numAggregators)
	aggResultChan := make(chan map[string]int, numAggregators)
	var aggWg sync.WaitGroup

	for i := 0; i < numAggregators; i++ {
		aggChans[i] = make(chan map[string]int, 100)
		aggWg.Add(1)
		go func(idx int) {
			defer aggWg.Done()
			localAgg := make(map[string]int)
			for localMap := range aggChans[idx] {
				for word, count := range localMap {
					localAgg[word] += count
				}
			}
			aggResultChan <- localAgg
		}(i)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range fileChan {
				localFreq := make(map[string]int)
				lines, err := processFileParallel(file, localFreq)
				if err != nil {
					fmt.Printf("[Worker %d] Error reading %s: %v\n", workerID, file, err)
					continue // Skip sending data for this file
				}

				fmt.Printf("[Worker %d] File: %s, Lines: %d\n", workerID, file, lines)

				// Only process successful reads
				for word, count := range localFreq {
					idx := int(hashWord(word)) % numAggregators
					aggChans[idx] <- map[string]int{word: count}
				}

				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}(i)
	}

	// Only one traversal here
	fileCount, walkErr := walkFiles(root, fileChan)
	if walkErr != nil {
		fmt.Printf("Error walking directory: %v\n", walkErr)
		return
	}
	fmt.Printf("Found %d text files. Using %d workers.\n", fileCount, numWorkers)

	wg.Wait()
	for _, ch := range aggChans {
		close(ch)
	}
	aggWg.Wait()

	aggStart := time.Now()
	wordFreq := make(map[string]int)
	for i := 0; i < numAggregators; i++ {
		local := <-aggResultChan
		for word, count := range local {
			wordFreq[word] += count
		}
	}
	close(aggResultChan)
	aggElapsed := time.Since(aggStart)
	fmt.Printf("Aggregation time: %s\n", aggElapsed)

	topStart := time.Now()
	printTopWords(wordFreq)
	topElapsed := time.Since(topStart)
	fmt.Printf("Top words calculation time: %s\n", topElapsed)

	totalElapsed := time.Since(walkStart)
	fmt.Printf("Total execution time: %s\n", totalElapsed)
}

// Returns the number of .txt files found and any error encountered
func walkFiles(root string, fileChan chan<- string) (int, error) {
	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			fileChan <- path
			count++
		}
		return nil
	})
	close(fileChan)
	return count, err
}

func processFileParallel(filename string, localFreq map[string]int) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Increase scanner buffer for large files
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 1024*1024) // 1MB buffer
	scanner.Buffer(buf, 1024*1024)

	lineCount := 0

	for scanner.Scan() {
		lineCount++
		for _, word := range strings.Fields(scanner.Text()) {
			normWord := normalizeWord(word)
			if normWord != "" {
				localFreq[normWord]++
			}
		}
	}
	return lineCount, scanner.Err()
}

// heap implementation for top K words
// if the size of the file is then this will help in performance
// time complexity O(N log K) instead of O(N log N) for sorting all words
type pair struct {
	word  string
	count int
}

type minHeap []pair

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(pair)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func printTopWords(wordFreq map[string]int) {
	const K = 10
	h := &minHeap{}
	heap.Init(h)
	for word, count := range wordFreq {
		heap.Push(h, pair{word, count})
		if h.Len() > K {
			heap.Pop(h)
		}
	}

	var result []pair
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(pair))
	}
	// Reverse for descending order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	fmt.Println("Top 10 words:")
	for _, p := range result {
		fmt.Printf("%s: %d\n", p.word, p.count)
	}
}

func normalizeWord(word string) string {
	var b strings.Builder
	b.Grow(len(word))
	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func hashWord(word string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(word))
	return h.Sum32()
}
