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
	"log"
	"runtime"
)

func filterOS(artifactDefinitions []ArtifactDefinition) []ArtifactDefinition {
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

func searchArtifacts(names []string, artifactDefinitions []ArtifactDefinition) []ArtifactDefinition {
	artifactDefinitionMap := map[string]ArtifactDefinition{}
	for _, artifactDefinition := range artifactDefinitions {
		artifactDefinitionMap[artifactDefinition.Name] = artifactDefinition
	}
	var artifactList []ArtifactDefinition
	for _, artifact := range getArtifacts(names, artifactDefinitionMap) {
		artifactList = append(artifactList, artifact)
	}
	return  artifactList
}

func getArtifacts(names []string, artifactDefinitions map[string]ArtifactDefinition) map[string]ArtifactDefinition {
	selected := map[string]ArtifactDefinition{}
	for _, name := range names {
		artifact, ok := artifactDefinitions[name]
		if !ok {
			log.Printf("Artifact Definition %s not found", name)
			continue
		}

		onlyGroup := true
		for _, source := range artifact.Sources {
			if source.Type == "ARTIFACT_GROUP" {
				for subName, subArtifact := range getArtifacts(source.Attributes.Names, artifactDefinitions){
					selected[subName] = subArtifact
				}
			} else {
				onlyGroup = false
			}
		}
		if !onlyGroup {
			selected[artifact.Name] = artifact
		}
	}

	return selected
}
