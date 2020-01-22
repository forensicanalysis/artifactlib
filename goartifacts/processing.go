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

// ProcessFiles takes a list of artifact definition files. Those files are decoded, validated, filtered and expanded.
func ProcessFiles(artifacts []string, filenames []string, collector ArtifactCollector) error {
	var artifactDefinitions []ArtifactDefinition

	// decode file
	for _, filename := range filenames {
		ads, _, err := DecodeFile(filename)
		if err != nil {
			return err
		}
		artifactDefinitions = append(artifactDefinitions, ads...)
	}

	return ProcessArtifacts(artifacts, artifactDefinitions, collector)
}

// ProcessArtifacts takes a list of artifact definitions. Those artifact definitions are filtered and expanded.
func ProcessArtifacts(artifacts []string, artifactDefinitions []ArtifactDefinition, collector ArtifactCollector) error {
	// select from entrypoint
	if artifacts != nil {
		artifactDefinitions = filterName(artifacts, artifactDefinitions)
	}

	// select supported os
	artifactDefinitions = filterOS(artifactDefinitions)

	// expand and glob
	Collect(artifactDefinitions, collector)

	return nil
}
