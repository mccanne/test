// Package test provides infrastructure for creating templated tests that may
// driven by the go test framework.  A test implements Interface and comprises
// a set of named test vectors forming a Suite.
package test

import (
	"errors"

	"github.com/stretchr/testify/require"
)

var ErrBadPhase = errors.New("phase number out of range")

type Params interface{}

type Vector interface {
	Params(id int) Params
}

func VectorParams(v Vector) []Params {
	var out []Params
	k := 0
	for {
		p := v.Params(k)
		if p == nil {
			return out
		}
		out = append(out, p)
		k++
	}
}

type Interface interface {
	Setup(int) error
	Vector() Vector
	Run(require.TestingT, Params)
	Teardown(int) error
	// true if this test shouldn't be run on a "short" pass
	IsLong() bool
}

type Long struct{}

func (p *Long) IsLong() bool {
	return true
}

type Short struct{}

func (p *Short) IsLong() bool {
	return false
}

type NopTeardown struct{}

func (p *NopTeardown) Teardown(phase int) error {
	return nil
}

type Array []Params

func (a *Array) Add(v interface{}) {
	*a = append(*a, v)
}

func (a Array) Params(id int) Params {
	if id >= 0 && id < len(a) {
		return a[id]
	}
	return nil
}
