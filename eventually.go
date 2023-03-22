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

// New creates a new Eventually with the given options. This can be useful if you want to reuse the same
// configuration for multiple functions. For example:
//
//	e := eventually.New(eventually.WithMaxAttempts(10))
//
// The returned Eventually has the following defaults unless otherwise specified:
//
//	Timeout:     10 seconds
//	Interval:    100 milliseconds
//	MaxAttempts: 0 (unlimited)
//
// If you don't need to reuse the same configuration, you can use the [Must] and [Should] functions directly.
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

// Must will keep retrying the given function f until the testing.TB passed to
// it does not fail or one of the following conditions is met:
//
//   - the timeout is reached
//   - the maximum number of attempts is reached
//
// If f does not succed, Must will halt the test calling t.Fatalf.
func (e *Eventually) Must(t testing.TB, f func(t testing.TB)) {
	t.Helper()

	r := e.retryableT(t)
	keepTrying(t, r, f, t.Fatalf)
}

// Should will keep retrying the given function f until the testing.TB passed to
// it does not fail or one of the following conditions is met:
//
//   - the timeout is reached
//   - the maximum number of attempts is reached
//
// If f does not succed, Should will fail the test calling t.Errorf.
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

// Option is a function that can be used to configure an Eventually.
type Option func(*Eventually)

// WithTimeout sets the timeout for an Eventually.
func WithTimeout(timeout time.Duration) Option {
	return func(e *Eventually) {
		e.timeout = timeout
	}
}

// WithInterval sets the interval Eventually will wait between attempts.
func WithInterval(interval time.Duration) Option {
	return func(e *Eventually) {
		e.interval = interval
	}
}

// WithMaxAttempts sets the maximum number of attempts an Eventually will make.
func WithMaxAttempts(attempts int) Option {
	return func(e *Eventually) {
		e.maxAttempts = attempts
	}
}

// Must will keep retrying the given function f until the testing.TB passed to
// it does not fail or one of the following conditions is met:
//
//   - the timeout is reached
//   - the maximum number of attempts is reached
//
// If f does not succed, Must will halt the test calling t.Fatalf.
// Must behaviour can be changed by passing options to it.
func Must(t testing.TB, f func(t testing.TB), options ...Option) {
	t.Helper()

	e := New(options...)
	e.Must(t, f)
}

// Should will keep retrying the given function f until the testing.TB passed to
// it does not fail or one of the following conditions is met:
//
//   - the timeout is reached
//   - the maximum number of attempts is reached
//
// If f does not succed, Should will fail the test calling t.Errorf.
// Should behaviour can be changed by passing options to it.
func Should(t testing.TB, f func(t testing.TB), options ...Option) {
	t.Helper()

	e := New(options...)
	e.Should(t, f)
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
