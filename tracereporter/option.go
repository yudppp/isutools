package tracereporter

import (
	"github.com/najeira/measure"
)

// Config .
type Config struct {
	serviceName string
	sortKey     string
}

// Option .
type Option func(*Config)

func defaults(cfg *Config) {
	cfg.sortKey = measure.Sum
}

// WithServiceName .
func WithServiceName(name string) Option {
	return func(cfg *Config) {
		cfg.serviceName = name
	}
}

// WithSortKey .
// sortKey is key | count | sum | min | max | avg | rate | p95
func WithSortKey(sortKey string) Option {
	return func(cfg *Config) {
		cfg.sortKey = sortKey
	}
}
