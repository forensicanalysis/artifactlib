// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package goartifacts

import (
	"github.com/forensicanalysis/fslib"
	"sync"
)

type ArtifactCollector interface {
	Resolve(parameter string) ([]string, error)
	Collect(name string, source Source)

	FS() fslib.FS
	Registry() fslib.FS
	AddPartitions() bool
}

func CollectAll(collector ArtifactCollector, artifactDefinitions []ArtifactDefinition) {
	var wg sync.WaitGroup
	for ax, artifactDefinition := range artifactDefinitions {
		wg.Add(1)
		go func(ax int, artifactDefinition ArtifactDefinition) {
			for _, source := range artifactDefinition.Sources {
				collector.Collect(artifactDefinition.Name, source)
			}
			wg.Done()
		}(ax, artifactDefinition)
	}
	wg.Wait()
}
