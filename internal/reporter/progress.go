package reporter

import (
	"fmt"
	"time"
)

type ProgressReporter struct {
	Total     int
	Processed int
	startTime time.Time
	mu        sync.Mutex
}

func NewProgressReporter(total int) *ProgressReporter {
	return &ProgressReporter{
		Total:     total,
		startTime: time.Now(),
	}
}

func (p *ProgressReporter) Update(count int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Processed += count
	
	elapsed := time.Since(p.startTime)
	rate := float64(p.Processed) / elapsed.Seconds()
	remaining := float64(p.Total-p.Processed) / rate
	
	fmt.Printf("\rProgress: %d/%d (%.2f/s, ETA: %.1fs)", 
		p.Processed, p.Total, rate, remaining)
}

func (p *ProgressReporter) Complete() {
	fmt.Printf("\nCompleted in %v\n", time.Since(p.startTime))
}
