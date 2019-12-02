package test

import (
	"fmt"
	"os"
)

// T implements the require.TestingT interface
type T struct {
	err error
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...)
}

func (t *T) FailNow() {
	fmt.Fprintln(os.Stderr, "failed test... exiting")
	if t.err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", t.err)
	}
	os.Exit(-1)
}
