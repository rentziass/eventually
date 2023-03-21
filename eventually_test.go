package eventually_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rentziass/eventually"
)

type test struct {
	*testing.T
	logs   []string
	failed bool
	halted bool
}

func (t *test) Fail() {
	t.failed = true
}

func (t *test) FailNow() {
	t.failed = true
	t.halted = true
}

func (t *test) Fatal(args ...interface{}) {
	t.Log(args...)
	t.FailNow()
}

func (t *test) Fatalf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.FailNow()
}

func (t *test) Error(args ...interface{}) {
	t.Log(args...)
	t.Fail()
}

func (t *test) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Fail()
}

func (t *test) Log(args ...any) {
	t.logs = append(t.logs, fmt.Sprintln(args...))
}

func (t *test) Logf(format string, args ...any) {
	t.logs = append(t.logs, fmt.Sprintf(format, args...))
}

func TestMust(t *testing.T) {
	t.Run("eventually succeeding", func(t *testing.T) {
		tt := &test{T: t}

		succeed := false

		eventually.Must(
			tt,
			func(t testing.TB) {
				if !succeed {
					t.Fail()
					succeed = true
				}
			},
			eventually.WithTimeout(100*time.Millisecond),
			eventually.WithInterval(1*time.Nanosecond),
		)

		if tt.failed || tt.halted {
			t.Error("test failed")
		}
	})

	t.Run("eventually failing", func(t *testing.T) {
		tt := &test{T: t}

		eventually.Must(
			tt,
			func(t testing.TB) {
				t.Fail()
			},
			eventually.WithTimeout(100*time.Millisecond),
			eventually.WithInterval(1*time.Nanosecond),
		)

		if !tt.failed {
			t.Error("test succeeded")
		}

		if !tt.halted {
			t.Error("test did not halt")
		}
	})

	t.Run("logs", func(t *testing.T) {
		tt := &test{T: t}

		eventually.Must(
			tt,
			func(t testing.TB) {
				t.Log("log")
			},
			eventually.WithMaxAttempts(1),
		)

		if len(tt.logs) != 1 {
			t.Fatalf("logs should contain 1 line, contained %d", len(tt.logs))
		}

		if tt.logs[0] != "log\n" {
			t.Errorf("logs should contain 'log', contained %q", tt.logs[0])
		}
	})
}

func TestShould(t *testing.T) {
	t.Run("eventually succeeding", func(t *testing.T) {
		tt := &test{T: t}

		succeed := false

		eventually.Should(
			tt,
			func(t testing.TB) {
				if !succeed {
					t.Fail()
					succeed = true
				}
			},
			eventually.WithTimeout(100*time.Millisecond),
			eventually.WithInterval(1*time.Nanosecond),
		)

		if tt.failed || tt.halted {
			t.Error("test failed")
		}
	})

	t.Run("eventually failing", func(t *testing.T) {
		tt := &test{T: t}

		eventually.Should(
			tt,
			func(t testing.TB) {
				t.Fail()
			},
			eventually.WithTimeout(100*time.Millisecond),
			eventually.WithInterval(1*time.Nanosecond),
		)

		if !tt.failed {
			t.Error("test succeeded")
		}

		if tt.halted {
			t.Error("Should should not halt")
		}
	})

	t.Run("logs", func(t *testing.T) {
		tt := &test{T: t}

		eventually.Should(
			tt,
			func(t testing.TB) {
				t.Log("log")
			},
			eventually.WithMaxAttempts(1),
		)

		if len(tt.logs) != 1 {
			t.Fatalf("logs should contain 1 line, contained %d", len(tt.logs))
		}

		if tt.logs[0] != "log\n" {
			t.Errorf("logs should contain 'log', contained %q", tt.logs[0])
		}
	})
}
