// Copyright (c) 2019 Siemens AG
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
	"fmt"
	"github.com/forensicanalysis/fslib"
)

// NamedSource wraps a Source to add the artifact name
type NamedSource struct {
	Source
	Name string
}

// ProcessFiles takes a list of artifact definition files. Those files are decoded, validated, filtered and expanded.
func ProcessFiles(artifacts []string, infs fslib.FS, addPartitions bool, filenames []string, resolver ParameterResolver) (artifactDefinitions []ArtifactDefinition, err error) {

	// decode file
	for _, filename := range filenames {
		ads, _, err := DecodeFile(filename)
		if err != nil {
			return artifactDefinitions, err
		}
		artifactDefinitions = append(artifactDefinitions, ads...)
	}

	// select from entrypoint
	if artifacts != nil {
		artifactDefinitions = filterName(artifacts, artifactDefinitions)
	}

	// select supported os
	artifactDefinitions = filterOS(artifactDefinitions)

	// expand and glob
	artifactDefinitions, err = Expand(infs, artifactDefinitions, addPartitions, resolver)

	return artifactDefinitions, err
}

// ParallelProcessArtifacts takes a list of artifact definitions. Those artifact definitions are filtered and expanded.
func ParallelProcessArtifacts(artifacts []string, infs fslib.FS, addPartitions bool, artifactDefinitions []ArtifactDefinition, resolver ParameterResolver) (<-chan NamedSource, int, error) {
	// select from entrypoint
	if artifacts != nil {
		artifactDefinitions = filterName(artifacts, artifactDefinitions)
	}

	// select supported os
	artifactDefinitions = filterOS(artifactDefinitions)

	sourceCount := 0
	for _, a := range artifactDefinitions {
		sourceCount += len(a.Sources)
	}

	sourceChannel := make(chan NamedSource, 100)
	// expand and glob
	go func() {
		ExpandChannel(sourceChannel, infs, artifactDefinitions, addPartitions, resolver)
		close(sourceChannel)
	}()
	return sourceChannel, sourceCount, nil
}

// ParallelProcessFiles takes a list of artifact definition files. Those files are decoded, validated, filtered and expanded.
func ParallelProcessFiles(artifacts []string, infs fslib.FS, addPartitions bool, filenames []string, resolver ParameterResolver) (<-chan NamedSource, int, error) {
	artifactDefinitionMap := map[string][]ArtifactDefinition{}
	var artifactDefinitions []ArtifactDefinition

	// decode file
	for _, filename := range filenames {
		ads, flaws, err := DecodeFile(filename)
		if err != nil {
			return nil, 0, err
		}
		if len(flaws) > 0 {
			fmt.Println(flaws)
		}
		artifactDefinitions = append(artifactDefinitions, ads...)
		artifactDefinitionMap[filename] = ads
	}

	// select from entrypoint
	if artifacts != nil {
		artifactDefinitions = filterName(artifacts, artifactDefinitions)
	}

	// select supported os
	artifactDefinitions = filterOS(artifactDefinitions)

	sourceCount := 0
	for _, a := range artifactDefinitions {
		sourceCount += len(a.Sources)
	}

	sourceChannel := make(chan NamedSource, 100)
	// expand and glob
	go func() {
		ExpandChannel(sourceChannel, infs, artifactDefinitions, addPartitions, resolver)
		close(sourceChannel)
	}()
	return sourceChannel, sourceCount, nil
}
