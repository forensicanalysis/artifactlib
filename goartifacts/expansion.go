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
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/filesystem/registryfs"
	"github.com/forensicanalysis/fslib/forensicfs/glob"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Expand performs parameter expansion and globbing on a list of artifact definitions.
func Expand(infs fslib.FS, artifactDefinitions []ArtifactDefinition, addPartitions bool) ([]ArtifactDefinition, error) {
	for ax, artifactDefinition := range artifactDefinitions {
		for sx, source := range artifactDefinition.Sources {
			artifactDefinitions[ax].Sources[sx] = expandSource(source, infs, addPartitions)
		}
	}
	return artifactDefinitions, nil
}

// ExpandChannel performs parameter expansion and globbing on a list of artifact definitions.
func ExpandChannel(sourceChannel chan<- NamedSource, infs fslib.FS, artifactDefinitions []ArtifactDefinition, addPartitions bool) {
	var wg sync.WaitGroup
	for ax, artifactDefinition := range artifactDefinitions {
		wg.Add(1)
		go func(ax int, artifactDefinition ArtifactDefinition) {
			for _, source := range artifactDefinition.Sources {
				sourceChannel <- NamedSource{expandSource(source, infs, addPartitions), artifactDefinition.Name}
			}
			wg.Done()
		}(ax, artifactDefinition)
	}
	wg.Wait()
}

func expandSource(source Source, infs fslib.FS, addPartitions bool) Source {
	replacer := strings.NewReplacer("\\", "/", "/", "\\")
	switch source.Type {
	case "FILE", "DIRECTORY", "PATH":
		// expand paths
		var expandedPaths []string
		for _, path := range source.Attributes.Paths {
			if source.Attributes.Separator == "\\" {
				path = strings.ReplaceAll(path, "\\", "/")
			}
			paths, err := expandPath(infs, path, addPartitions)
			if err != nil {
				log.Println(err)
				continue
			}
			expandedPaths = append(expandedPaths, paths...)
		}
		source.Attributes.Paths = expandedPaths
	case "REGISTRY_KEY":
		// expand keys
		var expandKeys []string
		for _, key := range source.Attributes.Keys {
			key = "/" + replacer.Replace(key)
			keys, err := expandKey(key)
			if err != nil {
				log.Println(err)
				continue
			}
			expandKeys = append(expandKeys, keys...)
		}
		source.Attributes.Keys = expandKeys
	case "REGISTRY_VALUE":
		// expand key value pairs
		var expandKeyValuePairs []KeyValuePair
		for _, keyValuePair := range source.Attributes.KeyValuePairs {
			key := "/" + replacer.Replace(keyValuePair.Key)
			keys, err := expandKey(key)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, expandKey := range keys {
				expandKeyValuePairs = append(expandKeyValuePairs, KeyValuePair{Key: expandKey, Value: keyValuePair.Value})
			}
		}
		source.Attributes.KeyValuePairs = expandKeyValuePairs
	}
	return source
}

func expandPath(fs fslib.FS, syspath string, addPartitions bool) ([]string, error) {
	// expand vars
	variablePaths := expandVar(syspath)
	if len(variablePaths) == 0 {
		return nil, nil
	}

	var forensicPaths []string
	for _, variablePath := range variablePaths {
		if addPartitions {
			var err error
			forensicPath, err := osfs.ToForensicPath(variablePath)
			if err != nil {
				return nil, err
			}
			forensicPaths = append(forensicPaths, forensicPath)
		} else {
			forensicPaths = append(forensicPaths, variablePath)
		}
	}

	// Test if variable path starts with e.g. C:/; need to be done after variable replacement
	isAbsPath, err := regexp.MatchString(`[a-zA-Z]:/`, variablePaths[0])
	if err != nil {
		return nil, err
	}

	var partitionPaths []string
	if runtime.GOOS == "windows" && addPartitions && !isAbsPath {
		partitions, err := listPartitions()
		if err != nil {
			return nil, err
		}
		for _, forensicPath := range forensicPaths {
			for _, partition := range partitions {
				partitionPaths = append(partitionPaths, fmt.Sprintf("/%s/%s", partition, forensicPath[3:]))
			}
		}
	} else {
		partitionPaths = forensicPaths
	}

	// unglob paths
	var unglobedPaths []string
	for _, expandedPath := range partitionPaths {
		expandedPath = strings.ReplaceAll(expandedPath, "{", `\{`)
		expandedPath = strings.ReplaceAll(expandedPath, "}", `\}`)
		unglobedPath, err := glob.Glob(fs, expandedPath)
		if err != nil {
			log.Println(err)
			continue
		}
		unglobedPaths = append(unglobedPaths, unglobedPath...)
	}

	return unglobedPaths, nil
}

func expandKey(path string) ([]string, error) {
	if runtime.GOOS == "windows" {
		return expandPath(registryfs.New(), path, false)
	}
	return []string{}, nil
}
