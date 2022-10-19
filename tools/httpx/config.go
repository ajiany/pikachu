package httpx

import (
	"time"
)

type Config struct {
	Host      string
	UrlPrefix string
	Timeout   time.Duration
	Trace     bool
}

func defaultCfg() *Config {
	return &Config{
		Timeout: 10 * time.Second,
	}
}

type configFn func(cfg *Config)

func WithHost(host string) configFn {
	return func(cfg *Config) {
		cfg.Host = host
	}
}

func WithUrlPrefix(prefix string) configFn {
	return func(cfg *Config) {
		cfg.UrlPrefix = prefix
	}
}

func WithTimeout(timeout time.Duration) configFn {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

func WithTrace(trace bool) configFn {
	return func(cfg *Config) {
		cfg.Trace = trace
	}
}
