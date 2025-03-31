package deleter

import (
	"sync"
	
	"github.com/yourusername/rmrf/internal/config"
	"github.com/yourusername/rmrf/internal/reporter"
)

type Deleter struct {
	config  *config.Options
	stats   *reporter.Stats
	mu      sync.Mutex
}

func New(opts ...config.Option) *Deleter {
	cfg := config.DefaultOptions
	for _, opt := range opts {
		opt(&cfg)
	}
	
	return &Deleter{
		config: &cfg,
		stats:  reporter.DefaultStats(),
	}
}

func (d *Deleter) Delete(path string) (*reporter.Stats, error) {
	if err := d.validatePath(path); err != nil {
		return nil, err
	}
	
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, d.config.MaxThreads)
	progress := reporter.NewProgressReporter(0) // Initialize with 0, will update during traversal

	wg.Add(1)
	go d.deleteRecursive(absPath, &wg, sem, progress)
	wg.Wait()
	progress.Complete()

	return d.stats, nil
}
