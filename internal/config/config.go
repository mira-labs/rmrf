package config

import "runtime"

type Options struct {
	MaxThreads    int
	DryRun        bool
	Interactive   bool
	Verbose       bool
	SkipSymlinks  bool
	DangerousPaths []string
}

type Option func(*Options)

func WithMaxThreads(n int) Option {
	return func(o *Options) {
		o.MaxThreads = n
	}
}

func WithDryRun(enabled bool) Option {
	return func(o *Options) {
		o.DryRun = enabled
	}
}

func WithInteractive(enabled bool) Option {
	return func(o *Options) {
		o.Interactive = enabled
	}
}

func WithVerbose(enabled bool) Option {
	return func(o *Options) {
		o.Verbose = enabled
	}
}
