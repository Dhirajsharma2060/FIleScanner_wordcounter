package main

import (
	"container/heap"
	"os"
	"strings"
	"testing"
)

func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Hello", "hello"},
		{"WORLD!", "world"},
		{"Go123", "go123"},
		{"", ""},
		{"!!!", ""},
		{"Go-Lang.", "golang"},
	}
	for _, tt := range tests {
		got := normalizeWord(tt.in)
		if got != tt.want {
			t.Errorf("normalizeWord(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestProcessFileParallel(t *testing.T) {
	content := "Hello world\nHello Go\nGo Go Go!"
	tmpfile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	localFreq := make(map[string]int)
	lines, err := processFileParallel(tmpfile.Name(), localFreq)
	if err != nil {
		t.Fatalf("processFileParallel error: %v", err)
	}
	if lines != 3 {
		t.Errorf("expected 3 lines, got %d", lines)
	}
	if localFreq["hello"] != 2 || localFreq["go"] != 4 || localFreq["world"] != 1 {
		t.Errorf("unexpected word counts: %+v", localFreq)
	}
}

func TestMinHeapTopK(t *testing.T) {
	wordFreq := map[string]int{
		"a": 10, "b": 20, "c": 30, "d": 40, "e": 50,
		"f": 60, "g": 70, "h": 80, "i": 90, "j": 100,
		"k": 110, "l": 120,
	}
	const K = 5
	h := &minHeap{}
	for word, count := range wordFreq {
		heap.Push(h, pair{word, count})
		if h.Len() > K {
			heap.Pop(h)
		}
	}
	if h.Len() != K {
		t.Errorf("expected heap size %d, got %d", K, h.Len())
	}
	// The smallest in the top K should be 80
	min := (*h)[0].count
	if min != 80 {
		t.Errorf("expected min count in heap to be 80, got %d", min)
	}
}

func TestWalkFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{"a.txt", "b.txt", "c.md"}
	for _, name := range files {
		os.WriteFile(dir+"/"+name, []byte("test"), 0644)
	}
	fileChan := make(chan string, 10)
	count, err := walkFiles(dir, fileChan)
	if err != nil {
		t.Fatalf("walkFiles error: %v", err)
	}
	got := []string{}
	for f := range fileChan {
		got = append(got, f)
	}
	if count != 2 {
		t.Errorf("expected 2 .txt files, got %d", count)
	}
	for _, f := range got {
		if !strings.HasSuffix(f, ".txt") {
			t.Errorf("unexpected file in channel: %s", f)
		}
	}
}
