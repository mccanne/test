package test

import (
	"fmt"
)

type Suite []Package

func (s Suite) LookupTest(name string, phase, id int) (*Package, Params, error) {
	pkg, err := s.LookupPackage(name)
	if err != nil {
		return nil, nil, err
	}
	if err := pkg.CheckPhase(phase); err != nil {
		return nil, nil, err
	}
	vector, err := pkg.PostVector(phase)
	params := vector.Params(id)
	if params == nil {
		return nil, nil, fmt.Errorf("test %d in package %s not found", id, name)
	}
	return pkg, params, nil
}

func (s Suite) LookupPackage(name string) (*Package, error) {
	for _, pkg := range s {
		if pkg.Name == name {
			return &pkg, nil
		}
	}
	return nil, fmt.Errorf("test package not found: %s", name)
}

func (s Suite) Run(name string, phase, id int) error {
	pkg, params, err := s.LookupTest(name, phase, id)
	if err != nil {
		return err
	}
	return pkg.Run(phase, id, params)
}

func (s Suite) RunAll(dataroot, name string) error {
	pkg, err := s.LookupPackage(name)
	if err != nil {
		return err
	}
	n := pkg.NumPhases()
	for phase := 0; phase < n; phase++ {
		vector, err := pkg.PostVector(phase)
		if err != nil {
			return err
		}
		if err := pkg.RunVector(VectorParams(vector), phase); err != nil {
			return err
		}
	}
	return nil
}
