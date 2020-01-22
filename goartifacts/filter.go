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
	"runtime"
	"strings"
)

func FilterOS(artifactDefinitions []ArtifactDefinition) []ArtifactDefinition {
	var selected []ArtifactDefinition
	for _, artifactDefinition := range artifactDefinitions {
		if isOSArtifactDefinition(runtime.GOOS, artifactDefinition.SupportedOs) {
			var sources []Source
			for _, source := range artifactDefinition.Sources {
				if isOSArtifactDefinition(runtime.GOOS, source.SupportedOs) {
					sources = append(sources, source)
				}
			}
			artifactDefinition.Sources = sources
			selected = append(selected, artifactDefinition)
		}
	}
	return selected
}

func FilterName(names []string, artifactDefinitions []ArtifactDefinition) []ArtifactDefinition {
	artifactDefinitionMap := map[string]ArtifactDefinition{}
	for _, artifactDefinition := range artifactDefinitions {
		artifactDefinitionMap[artifactDefinition.Name] = artifactDefinition
	}
	var artifactList []ArtifactDefinition
	for _, artifact := range expandArtifactGroup(names, artifactDefinitionMap) {
		artifactList = append(artifactList, artifact)
	}
	return  artifactList
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

func getSupportedOS(definition ArtifactDefinition, source Source) []string {
	if len(source.SupportedOs) > 0 {
		return source.SupportedOs
	} else if len(definition.SupportedOs) > 0 {
		return definition.SupportedOs
	}
	return listOSS()
}
