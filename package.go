package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type Package struct {
	Name   string
	Test   Interface
	Phases int
	// Baseline is nil for test packages that do not save their baseline
	// to the file system.  Otherwise, Baseline is a function that returns
	// an blank vector that can be used with json.Unmarshal so we get
	// proper type of json of objects instead of generic maps and arrays.
	Baseline func() Vector
	Dir      string
	Logger   *zap.Logger
	Terminal io.Writer
}

func (p *Package) log(format string, args ...interface{}) {
	w := p.Terminal
	if w != nil {
		fmt.Fprintf(w, format, args...)
	}
}

func (p *Package) PreVector() Vector {
	return p.Test.Vector()
}

func (p *Package) PostVector(phase int) (Vector, error) {
	if p.Baseline != nil {
		vector, err := p.LoadBaseline(p.baselinePath(phase))
		if err != nil {
			return nil, err
		}
		return vector, nil
	}
	return p.PreVector(), nil
}

func (p *Package) runTest(phase, id int, params Params) error {
	p.log("%s: running phase %d test %d\n", p.Name, phase, id)
	t := &T{}
	p.Test.Run(t, params)
	return t.err
}

func (p *Package) SetupPhase(targetPhase int) error {
	for phase := 0; phase < targetPhase; phase++ {
		p.log("%s: setup phase %d\n", p.Name, phase)
		if err := p.Test.Setup(phase); err != nil {
			return err
		}
		p.log("%s: teardown phase %d\n", p.Name, phase)
		if err := p.Test.Teardown(phase); err != nil {
			return err
		}
	}
	if err := p.Test.Setup(targetPhase); err != nil {
		return err
	}
	return nil
}

func (p *Package) CheckPhase(phase int) error {
	n := p.NumPhases()
	if phase < 0 {
		return fmt.Errorf("package \"%s\": phase %d not allowed", p.Name, phase)
	}
	if phase < n {
		return nil
	}
	s := "s"
	if n == 1 {
		s = ""
	}
	return fmt.Errorf("package \"%s\": phase %d not allowed as package has only %d phase%s", p.Name, phase, n, s)
}

func (p *Package) Run(phase, id int, params Params) error {
	if err := p.CheckPhase(phase); err != nil {
		return err
	}
	if err := p.SetupPhase(phase); err != nil {
		return err
	}
	if err := p.runTest(phase, id, params); err != nil {
		return err
	}
	return p.Test.Teardown(phase)
}

func (p *Package) NumPhases() int {
	n := p.Phases
	if n == 0 {
		n = 1
	}
	return n
}

func (p *Package) RunVector(vector []Params, phase int) error {
	if len(vector) == 0 {
		p.log("%s: warning: phase %d empty test vector... skipping\n", p.Name, phase)
		return nil
	}
	p.log("%s: setup phase %d\n", p.Name, phase)
	err := p.Test.Setup(phase)
	if err != nil {
		return err
	}
	for id, params := range vector {
		if err := p.runTest(phase, id, params); err != nil {
			return err
		}
	}
	p.log("%s: teardown phase %d\n", p.Name, phase)
	return p.Test.Teardown(phase)
}

func (p *Package) baselinePath(phase int) string {
	if p.Baseline == nil {
		return ""
	}
	name := fmt.Sprintf("%s.%d.json", p.Name, phase)
	return filepath.Join(p.Dir, "baseline", name)
}

func (p *Package) LoadBaseline(path string) (Vector, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v := p.Baseline()
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("%s: unmarshaling error: %s", path, err)
	}
	return v, nil
}

func (p *Package) SaveBaseline(path string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%s: marshaling error: %s", path, err)
	}
	return ioutil.WriteFile(path, data, 0600)
}

func (p *Package) UpdateBaseline() error {
	if p.Baseline == nil {
		return fmt.Errorf("package \"%s\" does not use a baseline", p.Name)
	}
	n := p.NumPhases()
	for phase := 0; phase < n; phase++ {
		vector := p.PreVector()
		if err := p.RunVector(VectorParams(vector), phase); err != nil {
			return err
		}
		path := p.baselinePath(phase)
		current, err := p.LoadBaseline(path)
		if err != nil {
			return err
		}
		if assert.ObjectsAreEqual(current, vector) {
			// XXX this doesn't seem to be working yet, but the
			// test comparisons work out for now.
			p.log("%s: baseline unchanged (file not modified)\n", path)
			return nil
		}
		if err := p.SaveBaseline(path, vector); err != nil {
			return err
		}
		p.log("%s: baseline updated\n", path)
	}
	return nil
}
