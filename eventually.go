package eventually

import (
	"testing"
	"time"
)

type failNowPanic struct{}

type retryableT struct {
	testing.TB

	failed bool

	timeout     time.Duration
	interval    time.Duration
	maxAttempts int
}

func (r *retryableT) Cleanup(func()) {
	// TODO: keep track of fuctions to run at the end of the current attempt
}

// TODO: should honor Skips too

func (r *retryableT) Error(args ...any) {
	r.Log(args...)
	r.Fail()
}

func (r *retryableT) Errorf(format string, args ...any) {
	r.Logf(format, args...)
	r.Fail()
}

func (r *retryableT) Fail() {
	r.failed = true
}

func (r *retryableT) FailNow() {
	r.failed = true
	panic(failNowPanic{})
}

func (r *retryableT) Failed() bool {
	return r.failed
}

func (r *retryableT) Fatal(args ...any) {
	r.Log(args...)
	r.FailNow()
}

func (r *retryableT) Fatalf(format string, args ...any) {
	r.Logf(format, args...)
	r.FailNow()
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
	keepTrying(t, f, t.Fatalf, options...)
}

func Should(t testing.TB, f func(t testing.TB), options ...Option) {
	keepTrying(t, f, t.Errorf, options...)
}

func keepTrying(t testing.TB, f func(t testing.TB), failf func(format string, args ...any), options ...Option) {
	retryable := &retryableT{
		TB: t,
	}

	for _, option := range options {
		option(retryable)
	}

	start := time.Now()
	attempts := 0

	for {
		if attempts >= retryable.maxAttempts && retryable.maxAttempts > 0 {
			// max attempts reached
			failf("eventually: max attempts reached")
			return
		}
		attempts++

		retryable.run(f)

		// test passed
		if !retryable.failed {
			break
		}

		if time.Since(start) >= retryable.timeout && retryable.timeout > 0 {
			// timeout reached
			failf("eventually: timeout reached")
			return
		}

		// test failed, wait for interval
		time.Sleep(retryable.interval)
	}
}

func (r *retryableT) run(f func(t testing.TB)) {
	r.failed = false

	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(failNowPanic); ok {
				return
			}

			panic(err)
		}
	}()

	f(r)
}
