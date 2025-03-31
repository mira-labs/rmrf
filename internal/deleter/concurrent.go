package deleter

import (
	"os"
	"path/filepath"
	"sync"
)

func (d *Deleter) deleteRecursive(path string, wg *sync.WaitGroup, sem chan struct{}, progress 
*reporter.ProgressReporter) {
	defer wg.Done()

	if err := d.makeDeletable(path); err != nil {
		d.stats.AddError(err)
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		d.stats.AddError(err)
		return
	}

	progress.Total += len(entries)
	var subWg sync.WaitGroup

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.Type()&os.ModeSymlink != 0 && d.config.SkipSymlinks {
			d.stats.AddError(fmt.Errorf("skipped symlink: %s", fullPath))
			continue
		}

		if entry.IsDir() {
			select {
			case sem <- struct{}{}:
				subWg.Add(1)
				go func(p string) {
					defer func() { <-sem }()
					d.deleteRecursive(p, &subWg, sem, progress)
				}(fullPath)
			default:
				d.deleteRecursive(fullPath, &subWg, sem, progress)
			}
		} else {
			d.processFile(fullPath)
			progress.Update(1)
		}
	}

	subWg.Wait()

	if !d.config.DryRun {
		if err := os.Remove(path); err != nil {
			d.stats.AddError(err)
		} else {
			d.stats.DirsDeleted++
		}
	}
}

func (d *Deleter) processFile(path string) {
	if d.config.DryRun {
		d.stats.FilesDeleted++
		return
	}

	if err := os.Chmod(path, 0600); err != nil {
		d.stats.AddError(err)
		return
	}

	if err := os.Remove(path); err != nil {
		d.stats.AddError(err)
	} else {
		d.stats.FilesDeleted++
	}
}
