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
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// ValidateFiles checks a list of files for various flaws.
func ValidateFiles(filenames []string) (flaws []Flaw, err error) {
	artifactDefinitionMap := map[string][]ArtifactDefinition{}

	// decode file
	for _, filename := range filenames {
		ads, typeflaw, err := DecodeFile(filename)
		if err != nil {
			return flaws, err
		}
		artifactDefinitionMap[filename] = ads
		flaws = append(flaws, typeflaw...)
	}

	// validate
	flaws = append(flaws, ValidateArtifactDefinitions(artifactDefinitionMap)...)
	return
}

// ProcessFiles takes a list of artifact definition files. Those files are decoded, validated, filtered and expanded.
func ProcessFiles(artifacts []string, infs fslib.FS, addPartitions bool, filenames []string) (artifactDefinitions []ArtifactDefinition, err error) {

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
		artifactDefinitions = searchArtifacts(artifacts, artifactDefinitions)
	}

	// select supported os
	artifactDefinitions = filterOS(artifactDefinitions)

	// expand and glob
	artifactDefinitions, err = Expand(infs, artifactDefinitions, addPartitions)

	return artifactDefinitions, err
}

// NamedSource wraps a Source to add the artifact name
type NamedSource struct {
	Source
	Name string
}

// ParallelProcessArtifacts takes a list of artifact definitions. Those artifact definitions are filtered and expanded.
func ParallelProcessArtifacts(artifacts []string, infs fslib.FS, addPartitions bool, artifactDefinitions []ArtifactDefinition) (<-chan NamedSource, int, error) {
	// select from entrypoint
	if artifacts != nil {
		artifactDefinitions = searchArtifacts(artifacts, artifactDefinitions)
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
		ExpandChannel(sourceChannel, infs, artifactDefinitions, addPartitions)
		close(sourceChannel)
	}()
	return sourceChannel, sourceCount, nil
}

// ParallelProcessFiles takes a list of artifact definition files. Those files are decoded, validated, filtered and expanded.
func ParallelProcessFiles(artifacts []string, infs fslib.FS, addPartitions bool, filenames []string) (<-chan NamedSource, int, error) {
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
		artifactDefinitions = searchArtifacts(artifacts, artifactDefinitions)
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
		ExpandChannel(sourceChannel, infs, artifactDefinitions, addPartitions)
		close(sourceChannel)
	}()
	return sourceChannel, sourceCount, nil
}

// DecodeFile takes a single artifact definition file to decode.
func DecodeFile(filename string) ([]ArtifactDefinition, []Flaw, error) {
	var artifactDefinitions []ArtifactDefinition
	var flaws []Flaw

	// open file
	f, err := os.Open(filename)
	if err != nil {
		return artifactDefinitions, flaws, err
	}
	defer f.Close()

	// decode file
	dec := NewDecoder(f)
	artifactDefinitions, err = dec.Decode()
	if err != nil {
		if typeerror, ok := err.(*yaml.TypeError); ok {
			// parsing error
			for _, typeerr := range typeerror.Errors {
				flaws = append(flaws, Flaw{Error, typeerr, "", filename})
			}
		} else {
			// bad error
			return artifactDefinitions, flaws, err
		}
	}

	return artifactDefinitions, flaws, nil
}

func isOSArtifactDefinition(os string, supportedOs []string) bool {
	if len(supportedOs) == 0 {
		return true
	}
	for _, supportedos := range supportedOs {
		if strings.EqualFold(supportedos, os) {
			return true
		}
	}
	return false
}
