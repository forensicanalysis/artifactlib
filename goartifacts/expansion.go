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
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/forensicfs/glob"
)

const windows = "windows"

// ExpandSource expands a single artifact definition source by expanding its
// paths or keys.
func ExpandSource(source Source, collector ArtifactCollector) Source {
	replacer := strings.NewReplacer("\\", "/", "/", "\\")
	switch source.Type {
	case SourceType.File, SourceType.Directory, SourceType.Path:
		// expand paths
		var expandedPaths []string
		for _, path := range source.Attributes.Paths {
			if source.Attributes.Separator == "\\" {
				path = strings.ReplaceAll(path, "\\", "/")
			}
			paths, err := expandPath(collector.FS(), path, collector.AddPartitions(), collector)
			if err != nil {
				log.Println(err)
				continue
			}
			expandedPaths = append(expandedPaths, paths...)
		}
		source.Attributes.Paths = expandedPaths
	case SourceType.RegistryKey:
		// expand keys
		var expandKeys []string
		for _, key := range source.Attributes.Keys {
			key = "/" + replacer.Replace(key)
			keys, err := expandKey(key, collector)
			if err != nil {
				log.Println(err)
				continue
			}
			expandKeys = append(expandKeys, keys...)
		}
		source.Attributes.Keys = expandKeys
	case SourceType.RegistryValue:
		// expand key value pairs
		var expandKeyValuePairs []KeyValuePair
		for _, keyValuePair := range source.Attributes.KeyValuePairs {
			key := "/" + replacer.Replace(keyValuePair.Key)
			keys, err := expandKey(key, collector)
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

func expandArtifactGroup(names []string, artifactDefinitions map[string]ArtifactDefinition) map[string]ArtifactDefinition {
	selected := map[string]ArtifactDefinition{}
	for _, name := range names {
		artifact, ok := artifactDefinitions[name]
		if !ok {
			log.Printf("Artifact Definition %s not found", name)
			continue
		}

		onlyGroup := true
		for _, source := range artifact.Sources {
			if source.Type == SourceType.ArtifactGroup {
				for subName, subArtifact := range expandArtifactGroup(source.Attributes.Names, artifactDefinitions) {
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

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func toForensicPath(name string, addPartitions bool) ([]string, error) { // nolint:gocyclo
	if runtime.GOOS != windows && name[0] != '/' {
		return nil, errors.New("path needs to be absolute")
	}

	if runtime.GOOS == windows {
		name = strings.ReplaceAll(name, `\`, "/")
		switch {
		case len(name) == 0:
			return []string{"/"}, nil
		case len(name) == 1:
			if name[0] == '/' {
				if addPartitions {
					root := &osfs.Root{}
					partitions, err := root.Readdirnames(0)
					if err != nil {
						return nil, err
					}
					var names []string
					for _, partition := range partitions {
						names = append(names, fmt.Sprintf("/%s", partition))
					}
					return names, nil
				}
				return []string{"/"}, nil
			} else if isLetter(name[0]) {
				return []string{"/" + name}, nil
			} else {
				return nil, fmt.Errorf("invalid path: %s", name)
			}
		case name[1] == ':':
			return []string{"/" + name[:1] + name[2:]}, nil
		case name[0] == '/' && isLetter(name[1]) && (len(name) == 2 || name[2] == '/'):
			return []string{name}, nil
		case addPartitions:
			root := &osfs.Root{}
			partitions, err := root.Readdirnames(0)
			if err != nil {
				return nil, err
			}
			var names []string
			for _, partition := range partitions {
				names = append(names, fmt.Sprintf("/%s%s", partition, name))
			}
			return names, nil
		default:
			return []string{name}, nil
		}
	}
	return []string{name}, nil
}

func expandPath(fs fslib.FS, syspath string, addPartitions bool, collector ArtifactCollector) ([]string, error) {
	// expand vars
	variablePaths, err := recursiveResolve(syspath, collector)
	if err != nil {
		return nil, err
	}
	if len(variablePaths) == 0 {
		return nil, nil
	}

	var partitionPaths []string
	for _, variablePath := range variablePaths {
		forensicPaths, err := toForensicPath(variablePath, addPartitions)
		if err != nil {
			return nil, err
		}
		partitionPaths = append(partitionPaths, forensicPaths...)
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

func expandKey(path string, collector ArtifactCollector) ([]string, error) {
	if runtime.GOOS == windows {
		return expandPath(collector.Registry(), path, false, collector)
	}
	return []string{}, nil
}

func recursiveResolve(s string, collector ArtifactCollector) ([]string, error) {
	var re = regexp.MustCompile(`%?%(.*?)%?%`)
	matches := re.FindAllStringSubmatch(s, -1)

	if len(matches) > 0 {
		var replacedParameters []string
		for _, match := range matches {
			resolves, err := collector.Resolve(match[1])
			if err != nil {
				return nil, err
			}

			replacedParameters = append(replacedParameters, replaces(re, s, resolves)...)
		}
		var results []string
		for _, result := range replacedParameters {
			childResults, err := recursiveResolve(result, collector)
			if err != nil {
				return nil, err
			}
			results = append(results, childResults...)
		}
		return results, nil
	}
	return []string{s}, nil
}

func replaces(regex *regexp.Regexp, s string, news []string) (ss []string) {
	for _, newString := range news {
		ss = append(ss, regex.ReplaceAllString(s, newString))
	}
	return
}
