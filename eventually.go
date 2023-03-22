package eventually

import (
	"testing"
	"time"
)

type Eventually struct {
	timeout     time.Duration
	interval    time.Duration
	maxAttempts int
}

func New(options ...Option) *Eventually {
	e := &Eventually{
		timeout:     10 * time.Second,
		interval:    100 * time.Millisecond,
		maxAttempts: 0,
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func (e *Eventually) Must(t testing.TB, f func(t testing.TB)) {
	t.Helper()

	r := e.retryableT(t)
	keepTrying(t, r, f, t.Fatalf)
}

func (e *Eventually) Should(t testing.TB, f func(t testing.TB)) {
	t.Helper()

	r := e.retryableT(t)
	keepTrying(t, r, f, t.Errorf)
}

func (e *Eventually) retryableT(t testing.TB) *retryableT {
	return &retryableT{
		TB:          t,
		timeout:     e.timeout,
		interval:    e.interval,
		maxAttempts: e.maxAttempts,
	}
}

type failNowPanic struct{}

type retryableT struct {
	testing.TB

	failed bool

	timeout     time.Duration
	interval    time.Duration
	maxAttempts int
}

func (r *retryableT) Error(args ...any) {
	r.TB.Helper()
	r.Log(args...)
	r.Fail()
}

func (r *retryableT) Errorf(format string, args ...any) {
	r.TB.Helper()
	r.Logf(format, args...)
	r.Fail()
}

func (r *retryableT) Fail() {
	r.TB.Helper()
	r.failed = true
}

func (r *retryableT) FailNow() {
	r.TB.Helper()
	r.failed = true
	panic(failNowPanic{})
}

func (r *retryableT) Failed() bool {
	return r.failed
}

func (r *retryableT) Fatal(args ...any) {
	r.TB.Helper()
	r.Log(args...)
	r.FailNow()
}

func (r *retryableT) Fatalf(format string, args ...any) {
	r.TB.Helper()
	r.Logf(format, args...)
	r.FailNow()
}

type Option func(*Eventually)

func WithTimeout(timeout time.Duration) Option {
	return func(e *Eventually) {
		e.timeout = timeout
	}
}

func WithInterval(interval time.Duration) Option {
	return func(e *Eventually) {
		e.interval = interval
	}
}

func WithMaxAttempts(attempts int) Option {
	return func(e *Eventually) {
		e.maxAttempts = attempts
	}
}

func Must(t testing.TB, f func(t testing.TB), options ...Option) {
	t.Helper()

	e := New(options...)
	e.Must(t, f)
}

func Should(t testing.TB, f func(t testing.TB), options ...Option) {
	t.Helper()

	e := New(options...)
	e.Should(t, f)
}

func keepTrying(t testing.TB, retryable *retryableT, f func(t testing.TB), failf func(format string, args ...any)) {
	t.Helper()

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
