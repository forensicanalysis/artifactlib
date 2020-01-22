package goartifacts

import (
	"errors"
	"github.com/forensicanalysis/fslib"
)

type TestCollector struct {
	fs fslib.FS
	Collected map[string][]Source
}

func (r *TestCollector) Collect(name string, source Source) {
	if r.Collected == nil {
		r.Collected =  map[string][]Source{}
	}
	r.Collected[name] = append(r.Collected[name], source)
}

func (r *TestCollector) FS() fslib.FS {
	return r.fs
}

func (r *TestCollector) Registry() fslib.FS {
	return r.fs
}

func (r *TestCollector) AddPartitions() bool {
	return false
}

func (r *TestCollector) Resolve(s string) ([]string, error) {
	switch s {
	case "foo":
		return []string{"xxx", "yyy"}, nil
	case "faz":
		return []string{"%foo%"}, nil
	}
	return nil, errors.New("could not resolve")
}
