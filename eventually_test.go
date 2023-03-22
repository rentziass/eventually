package eventually_test

import (
	"testing"
	"time"

	"github.com/rentziass/eventually"
)

func TestEventually_Must(t *testing.T) {
	t.Run("eventually succeeding", func(t *testing.T) {
		tt := &test{T: t}

		succeed := false

		e := eventually.New()
		e.Must(tt, func(t testing.TB) {
			if !succeed {
				t.Fail()
				succeed = true
			}
		})

		if tt.failed || tt.halted {
			t.Error("test failed")
		}
	})

	t.Run("eventually failing", func(t *testing.T) {
		tt := &test{T: t}

		e := eventually.New(
			eventually.WithMaxAttempts(3),
			eventually.WithInterval(0),
		)
		e.Must(tt, func(t testing.TB) {
			t.Fail()
		})

		if !tt.failed {
			t.Error("test succeeded")
		}

		if !tt.halted {
			t.Error("test did not halt")
		}
	})

	t.Run("logs", func(t *testing.T) {
		tt := &test{T: t}

		e := eventually.New(eventually.WithMaxAttempts(1))
		e.Must(tt, func(t testing.TB) {
			t.Log("hello")
			t.Log("world")
		})

		if len(tt.logs) != 2 {
			t.Fatalf("logs should contain 2 line, contained %d", len(tt.logs))
		}

		if tt.logs[0] != "hello\n" {
			t.Error("unexpected log")
		}

		if tt.logs[1] != "world\n" {
			t.Error("unexpected log")
		}
	})
}

func TestEventually_Should(t *testing.T) {
	t.Run("eventually succeeding", func(t *testing.T) {
		tt := &test{T: t}

		succeed := false

		e := eventually.New()
		e.Should(tt, func(t testing.TB) {
			if !succeed {
				t.Fail()
				succeed = true
			}
		})

		if tt.failed || tt.halted {
			t.Error("test failed")
		}
	})

	t.Run("eventually failing", func(t *testing.T) {
		tt := &test{T: t}

		e := eventually.New(
			eventually.WithMaxAttempts(3),
			eventually.WithInterval(0),
		)
		e.Should(tt, func(t testing.TB) {
			t.Fail()
		})

		if !tt.failed {
			t.Error("test succeeded")
		}

		if tt.halted {
			t.Error("test halted")
		}
	})

	t.Run("logs", func(t *testing.T) {
		tt := &test{T: t}

		e := eventually.New(eventually.WithMaxAttempts(1))
		e.Should(tt, func(t testing.TB) {
			t.Log("hello")
			t.Log("world")
		})

		if len(tt.logs) != 2 {
			t.Fatalf("logs should contain 2 line, contained %d", len(tt.logs))
		}

		if tt.logs[0] != "hello\n" {
			t.Error("unexpected log")
		}

		if tt.logs[1] != "world\n" {
			t.Error("unexpected log")
		}
	})

	t.Run("multiple uses", func(t *testing.T) {
		tt := &test{T: t}

		e := eventually.New(eventually.WithMaxAttempts(3), eventually.WithInterval(0))

		succeed := false
		e.Should(tt, func(t testing.TB) {
			if !succeed {
				t.Fail()
				t.Log("should")
				succeed = true
			}
		})

		e.Must(tt, func(t testing.TB) {
			t.Log("must")
		})

		if tt.failed || tt.halted {
			t.Error("test failed")
		}

		if len(tt.logs) != 2 {
			t.Fatalf("logs should contain 2 line, contained %d", len(tt.logs))
		}

		if tt.logs[0] != "should\n" {
			t.Error("unexpected log")
		}

		if tt.logs[1] != "must\n" {
			t.Error("unexpected log")
		}
	})
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
