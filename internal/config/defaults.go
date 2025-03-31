package config

var DefaultOptions = Options{
	MaxThreads:    runtime.NumCPU(),
	DryRun:        false,
	Interactive:   false,
	Verbose:       false,
	SkipSymlinks:  true,
	DangerousPaths: []string{"/", "/etc", "/usr", "/bin", "/sbin"},
}
