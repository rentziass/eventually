package eventually_test

import (
	"fmt"
	"testing"
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
