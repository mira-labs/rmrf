// Package main provides a fast, multi-threaded alternative to `rm -rf` using Go's goroutines.
// It efficiently deletes directories and their contents in parallel, optimizing for system resources.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// deleteRecursive removes all files and subdirectories inside the given directory concurrently.
// It uses a semaphore to limit the number of concurrent goroutines, preventing excessive CPU usage.
//
// Parameters:
//   - path: The directory path to delete.
//   - wg: A pointer to a sync.WaitGroup to track goroutine completion.
//   - semaphore: A buffered channel controlling concurrency limits.
func deleteRecursive(path string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done() // Mark this goroutine as completed.

	// Read all entries in the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", path, err)
		return
	}

	var subWg sync.WaitGroup // Sub-waitgroup to manage subdirectory deletions.

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// Acquire a semaphore slot before spawning a new goroutine.
			semaphore <- struct{}{}
			subWg.Add(1)
			go func(p string) {
				defer func() { <-semaphore }() // Release slot when done.
				deleteRecursive(p, &subWg, semaphore)
			}(fullPath)
		} else {
			// Delete the file
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Failed to delete file %s: %v\n", fullPath, err)
			}
		}
	}

	// Wait for all subdirectories to be deleted.
	subWg.Wait()

	// Finally, remove the now-empty directory itself.
	err = os.Remove(path)
	if err != nil {
		fmt.Printf("Failed to remove directory %s: %v\n", path, err)
	}
}

// asyncRmRF initiates a multi-threaded recursive directory deletion process.
// It automatically determines the optimal level of concurrency based on available CPU cores.
//
// Parameters:
//   - path: The directory to delete.
//
// Usage:
//
//	asyncRmRF("node_modules") // Deletes the "node_modules" folder asynchronously.
func asyncRmRF(path string) {
	var wg sync.WaitGroup
	maxThreads := runtime.NumCPU() // Use available CPU cores for concurrency.
	semaphore := make(chan struct{}, maxThreads)

	wg.Add(1)
	go deleteRecursive(path, &wg, semaphore)

	wg.Wait() // Wait for all deletion tasks to complete.
	fmt.Println("Deletion completed for", path)
}

// main handles command-line arguments and starts the deletion process.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<directory>")
		os.Exit(1)
	}

	asyncRmRF(os.Args[1])
}
