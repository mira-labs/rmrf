package main

import (
	"fmt"
	"os"
	
	"github.com/yourusername/rmrf/internal/deleter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	del := deleter.New(
		deleter.WithMaxThreads(8),
		deleter.WithDryRun(false),
	)

	stats, err := del.Delete(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDeletion complete:\n")
	fmt.Printf("- Files: %d\n", stats.FilesDeleted)
	fmt.Printf("- Directories: %d\n", stats.DirsDeleted)
	
	if len(stats.Errors) > 0 {
		fmt.Printf("\nEncountered %d errors:\n", len(stats.Errors))
		for _, err := range stats.Errors {
			fmt.Printf("  - %v\n", err)
		}
		os.Exit(1)
	}
}
