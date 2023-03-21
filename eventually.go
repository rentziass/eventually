package eventually

import (
	"testing"
	"time"
)

type retryableT struct {
	testing.TB

	failed bool

	timeout     time.Duration
	interval    time.Duration
	maxAttempts int
}

func (r *retryableT) Cleanup(func()) {
	// keep track of fuctions to run at the end of the current attempt
}

func (r *retryableT) Error(args ...any) {
	// write to ttb.Log
	// set r failed to true
}

func (r *retryableT) Errorf(format string, args ...any) {
	// write to ttb.Logf
	// set r failed to true
}

func (r *retryableT) Fail() {
	// set r failed to true
}

func (r *retryableT) FailNow() {
	// panic with a special error type
}

func (r *retryableT) Failed() bool {
	return r.failed
}

func (r *retryableT) Fatal(args ...any) {
	// write to ttb.Log
	// panic with a special error type
}

func (r *retryableT) Fatalf(format string, args ...any) {
	// write to ttb.Logf
	// panic with a special error type
}

type Option func(*retryableT)

func WithTimeout(timeout time.Duration) Option {
	return func(r *retryableT) {
		r.timeout = timeout
	}
}

func WithInterval(interval time.Duration) Option {
	return func(r *retryableT) {
		r.interval = interval
	}
}

func WithMaxAttempts(attempts int) Option {
	return func(r *retryableT) {
		r.maxAttempts = attempts
	}
}

func Must(t testing.TB, f func(t testing.TB), options ...Option) {
}

func Should(t testing.TB, f func(t testing.TB), options ...Option) {
}
