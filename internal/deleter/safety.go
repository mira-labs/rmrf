package deleter

import (
	"errors"
	"os"
)

var (
	ErrDangerousPath = errors.New("dangerous path specified")
	ErrNotExist      = errors.New("path does not exist")
)

func (d *Deleter) validatePath(path string) error {
	for _, dangerous := range d.config.DangerousPaths {
		if path == dangerous {
			return ErrDangerousPath
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrNotExist
	}

	return nil
}

func (d *Deleter) makeDeletable(path string) error {
	if d.config.DryRun {
		return nil
	}
	return os.Chmod(path, 0700)
}
