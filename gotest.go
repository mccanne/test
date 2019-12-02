package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type testingT struct {
	*testing.T
	cmd string
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	hdr := "=== RUN THIS TO REPRODUCE FAILED TEST ==="
	format += fmt.Sprintf("\n%s\n%s\n%s\n\n", hdr, t.cmd, hdr)
	t.T.Errorf(format, args...)
}

type FailCommands struct {
	Phase func(string, int) string
	Test  func(string, int, int) string
}

func RunGoTests(t *testing.T, dataroot string, fail FailCommands, packages []Package) {
	for _, pkg := range packages {
		pkg.Dir = dataroot
		if testing.Short() {
			if pkg.Test.IsLong() {
				t.Skip(fmt.Sprintf("package %s: skipping big test in short test run", pkg.Name))
			}
		}
		nphase := pkg.Phases
		if nphase == 0 {
			nphase = 1
		}
		for phase := 0; phase < nphase; phase++ {
			cmd := fail.Phase(pkg.Name, phase)
			T := &testingT{t, cmd}
			err := pkg.Test.Setup(phase)
			require.NoError(T, err)
			vector, err := pkg.PostVector(phase)
			require.NoError(T, err)
			name := fmt.Sprintf("%s:%d", pkg.Name, phase)
			t.Run(name, func(t *testing.T) {
				id := 0
				for {
					params := vector.Params(id)
					if params == nil {
						break
					}
					name := strconv.Itoa(id)
					cmd := fail.Test(pkg.Name, phase, id)
					t.Run(name, func(t *testing.T) {
						pkg.Test.Run(&testingT{t, cmd}, params)
					})
					id++
				}
			})
			err = pkg.Test.Teardown(phase)
			require.NoError(T, err)
		}
	}
}
