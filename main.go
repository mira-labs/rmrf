// Package main provides a fast, multi-threaded alternative to `rm -rf` with enhanced safety features.
//
// Key Features:
//   - Concurrent deletion using goroutines with semaphore-based throttling
//   - Configurable maximum thread usage (defaults to CPU cores)
//   - Comprehensive error handling and statistics collection
//   - Permission management (attempts chmod before deletion)
//   - Symlink protection (skips rather than follows)
//   - Version information embedded in builds
//
// Usage:
//
//	./rmrf <directory>
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// version will be set during build via ldflags (see Makefile)
var version = "development"

// DeleteStats tracks the results of a deletion operation for reporting
//
// Fields:
//   - FilesDeleted: Count of successfully deleted files
//   - DirsDeleted: Count of successfully deleted directories
//   - Errors: Count of errors encountered during deletion
type DeleteStats struct {
	FilesDeleted int
	DirsDeleted  int
	Errors       int
}

// deleteRecursive recursively deletes directory contents with concurrency control
//
// Parameters:
//   - path: Absolute path to delete
//   - wg: WaitGroup for tracking completion
//   - semaphore: Channel for limiting concurrent operations
//   - stats: Pointer to DeleteStats for collecting metrics
//   - mutex: Mutex for thread-safe stats modification
//
// Operation Flow:
//  1. Attempts to make directory writable
//  2. Reads all directory entries
//  3. Processes files and subdirectories concurrently
//  4. Waits for all suboperations to complete
//  5. Removes the now-empty directory
func deleteRecursive(path string, wg *sync.WaitGroup, semaphore chan struct{}, stats *DeleteStats, mutex *sync.Mutex) {
	defer wg.Done()

	// Try to make directory writable (log failure but continue)
	if err := os.Chmod(path, 0700); err != nil {
		log.Printf("warning: couldn't modify permissions for %s: %v", path, err)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		recordError(stats, mutex, fmt.Errorf("error reading directory %s: %w", path, err))
		return
	}

	var subWg sync.WaitGroup

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		// Skip symlinks to prevent accidental traversal
		if entry.Type()&os.ModeSymlink != 0 {
			log.Printf("skipping symlink: %s", fullPath)
			continue
		}

		if entry.IsDir() {
			select {
			case semaphore <- struct{}{}: // Acquire semaphore if possible
				subWg.Add(1)
				go func(p string) {
					defer func() { <-semaphore }() // Release when done
					deleteRecursive(p, &subWg, semaphore, stats, mutex)
				}(fullPath)
			default:
				// If semaphore full, process sequentially
				deleteRecursive(fullPath, &subWg, semaphore, stats, mutex)
			}
		} else {
			// Make file writable before deletion attempt
			if err := os.Chmod(fullPath, 0600); err != nil {
				recordError(stats, mutex, fmt.Errorf("couldn't modify file permissions %s: %w", fullPath, err))
			}

			if err := os.Remove(fullPath); err != nil {
				recordError(stats, mutex, fmt.Errorf("failed to delete file %s: %w", fullPath, err))
			} else {
				mutex.Lock()
				stats.FilesDeleted++
				mutex.Unlock()
			}
		}
	}

	subWg.Wait()

	// Remove the now-empty directory
	if err := os.Remove(path); err != nil {
		recordError(stats, mutex, fmt.Errorf("failed to remove directory %s: %w", path, err))
	} else {
		mutex.Lock()
		stats.DirsDeleted++
		mutex.Unlock()
	}
}

// recordError safely increments the error counter and logs the error
//
// Parameters:
//   - stats: DeleteStats instance to modify
//   - mutex: Mutex for thread-safe access
//   - err: Error to record
func recordError(stats *DeleteStats, mutex *sync.Mutex, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	stats.Errors++
	log.Println(err)
}

// asyncRmRF performs concurrent recursive directory deletion with safety checks
//
// Parameters:
//   - path: Target directory to delete
//
// Returns:
//   - *DeleteStats: Operation statistics
//   - error: Initialization errors (nil if operation started)
//
// Safety Mechanisms:
//   - Rejects empty, root, and current directory paths
//   - Resolves to absolute path before operation
//   - Limits maximum concurrent operations
func asyncRmRF(path string) (*DeleteStats, error) {
	// Basic safety checks
	if path == "" || path == "/" || path == "." {
		return nil, fmt.Errorf("dangerous path specified")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	var wg sync.WaitGroup
	var stats DeleteStats
	var mutex sync.Mutex

	maxThreads := runtime.NumCPU()
	semaphore := make(chan struct{}, maxThreads)

	wg.Add(1)
	go deleteRecursive(absPath, &wg, semaphore, &stats, &mutex)

	wg.Wait()
	return &stats, nil
}

// main is the CLI entry point that handles argument parsing and output
//
// Usage:
//
//	./rmrf <directory>
//
// Exit Codes:
//
//	0 - Success
//	1 - Usage error
//	2 - Deletion completed with errors
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Concurrent rm -rf tool v%s\n", version)
		fmt.Printf("Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	stats, err := asyncRmRF(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Deletion completed for %s\n", os.Args[1])
	fmt.Printf("Results: %d files, %d directories deleted, %d errors\n",
		stats.FilesDeleted, stats.DirsDeleted, stats.Errors)

	if stats.Errors > 0 {
		os.Exit(2)
	}
}
