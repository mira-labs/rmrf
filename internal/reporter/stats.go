package reporter

import (
	"encoding/json"
	"sync"
)

type Stats struct {
	FilesDeleted int      `json:"filesDeleted"`
	DirsDeleted  int      `json:"dirsDeleted"`
	Errors       []error  `json:"-"`
	mu           sync.Mutex
}

func DefaultStats() *Stats {
	return &Stats{
		Errors: make([]error, 0),
	}
}

func (s *Stats) AddError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Errors = append(s.Errors, err)
}

func (s *Stats) JSON() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, _ := json.Marshal(s)
	return string(data)
}
