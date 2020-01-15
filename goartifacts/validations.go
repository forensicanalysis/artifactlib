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
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/looplab/tarjan"
)

// global

func (r *validator) validateNameUnique(artifactDefinitions []ArtifactDefinition) {
	var knownNames = map[string]bool{}
	for _, artifactDefinition := range artifactDefinitions {
		if _, ok := knownNames[artifactDefinition.Name]; ok {
			r.addWarning("", artifactDefinition.Name, "Duplicate artifact name %s", artifactDefinition.Name)
		} else {
			knownNames[artifactDefinition.Name] = true
		}
	}
}

func (r *validator) validateRegistryKeyUnique(artifactDefinitions []ArtifactDefinition) {
	var knownKeys = map[string]bool{}
	for _, artifactDefinition := range artifactDefinitions {
		for _, source := range artifactDefinition.Sources {
			for _, key := range source.Attributes.Keys {
				if _, ok := knownKeys[key]; ok {
					r.addWarning("", artifactDefinition.Name, "Duplicate registry key %s", key)
				} else {
					knownKeys[key] = true
				}
			}
		}
	}
}

func (r *validator) validateRegistryValueUnique(artifactDefinitions []ArtifactDefinition) {
	var knownKeys = map[string]bool{}
	for _, artifactDefinition := range artifactDefinitions {
		for _, source := range artifactDefinition.Sources {
			for _, keyvalue := range source.Attributes.KeyValuePairs {
				if _, ok := knownKeys[keyvalue.Key+"/"+keyvalue.Value]; ok {
					r.addWarning("", artifactDefinition.Name, "Duplicate registry value %s %s", keyvalue.Key, keyvalue.Value)
				} else {
					knownKeys[keyvalue.Key+"/"+keyvalue.Value] = true
				}
			}
		}
	}
}

func (r *validator) validateNoCycles(artifactDefinitions []ArtifactDefinition) {
	graph := make(map[interface{}][]interface{})
	for _, artifactDefinition := range artifactDefinitions {
		for _, source := range artifactDefinition.Sources {
			if source.Type == "ARTIFACT_GROUP" {
				graph[artifactDefinition.Name] = []interface{}{}
				for _, name := range source.Attributes.Names {
					if name == artifactDefinition.Name {
						r.addError("", artifactDefinition.Name, "Artifact group references itself")
					}
					graph[artifactDefinition.Name] = append(graph[artifactDefinition.Name], name)
				}
			}
		}
	}

	output := tarjan.Connections(graph)
	for _, subgraph := range output {
		if len(subgraph) > 1 {
			var sortedSubgraph []string
			for _, subgraphitem := range subgraph {
				sortedSubgraph = append(sortedSubgraph, subgraphitem.(string))
			}
			sort.Strings(sortedSubgraph)
			r.addError("", "", "Cyclic artifact group: %s", sortedSubgraph)
		}
	}
}

func (r *validator) validateGroupMemberExist(artifactDefinitions []ArtifactDefinition) {
	var knownNames = map[string]bool{}
	for _, artifactDefinition := range artifactDefinitions {
		knownNames[artifactDefinition.Name] = true
	}

	for _, artifactDefinition := range artifactDefinitions {
		for _, source := range artifactDefinition.Sources {
			for _, member := range source.Attributes.Names {
				if _, ok := knownNames[member]; !ok {
					r.addError("", artifactDefinition.Name, "Unknown name %s in %s", member, artifactDefinition.Name)
				}
			}
		}
	}
}

func (r *validator) validateParametersProvided(artifactDefinitions []ArtifactDefinition) {
	var knownProvides = map[string]string{}
	for _, artifactDefinition := range artifactDefinitions {
		for _, provide := range artifactDefinition.Provides {
			knownProvides[provide] = artifactDefinition.Name
		}
	}

	/* for parameter := range knowledgeBase {
		if val, ok := knownProvides[parameter]; !ok {
			r.addInfo("", val, "Parameter %s is not provided", parameter)
		}
	}*/
}

// file

func (r *validator) validateNamePrefix(filename string, artifactDefinition ArtifactDefinition) {
	prefix := ""
	switch {
	case strings.HasPrefix(filepath.Base(filename), "windows"):
		prefix = "Windows"
	case strings.HasPrefix(filepath.Base(filename), "linux"):
		prefix = "Linux"
	case strings.HasPrefix(filepath.Base(filename), "macos"):
		prefix = "MacOS"
	}
	if !strings.HasPrefix(artifactDefinition.Name, prefix) {
		r.addCommon(filename, artifactDefinition.Name, "Artifact name should start with %s", prefix)
	}
}

func (r *validator) validateOSSpecific(filename string, artifactDefinition ArtifactDefinition) {
	os := ""
	if strings.HasPrefix(filepath.Base(filename), "windows") {
		os = supportedOS.Windows
	} else if strings.HasPrefix(filepath.Base(filename), "linux") {
		os = supportedOS.Linux
	} else if strings.HasPrefix(filepath.Base(filename), "macos") {
		os = supportedOS.Darwin
	}
	if os == "" {
		return
	}

	for _, supportedOs := range artifactDefinition.SupportedOs {
		if supportedOs != os {
			r.addInfo(filename, artifactDefinition.Name, "File should only contain %s artifact definitions", os)
		}
	}
	for _, source := range artifactDefinition.Sources {
		for _, supportedOs := range source.SupportedOs {
			if supportedOs != os {
				r.addInfo(filename, artifactDefinition.Name, "File should only contain %s artifact definitions", os)
			}
		}
	}
}

// artifact

func (r *validator) validateNameCase(filename string, artifactDefinition ArtifactDefinition) {
	if len(artifactDefinition.Name) < 2 {
		r.addError(filename, artifactDefinition.Name, "Artifact names be longer than 2 characters")
		return
	}
	if strings.ToUpper(artifactDefinition.Name[:1]) != artifactDefinition.Name[:1] {
		r.addInfo(filename, artifactDefinition.Name, "Artifact names should be CamelCase")
	}
	if strings.ContainsAny(artifactDefinition.Name, " \t") {
		r.addInfo(filename, artifactDefinition.Name, "Artifact names should not contain whitespace")
	}
}

func (r *validator) validateNameTypeSuffix(filename string, artifactDefinition ArtifactDefinition) {
	if len(artifactDefinition.Sources) == 0 {
		r.addError(filename, artifactDefinition.Name, "Artifact has no sources")
		return
	}
	currentSourceType := artifactDefinition.Sources[0].Type
	for _, source := range artifactDefinition.Sources {
		if source.Type != currentSourceType {
			return
		}
	}

	endings := map[string][]string{
		sourceType.Command:       {"Command", "Commands"},
		sourceType.Directory:     {"Directory", "Directories"},
		sourceType.File:          {"File", "Files"},
		sourceType.Path:          {"Path", "Paths"},
		sourceType.RegistryKey:   {"RegistryKey", "RegistryKeys"},
		sourceType.RegistryValue: {"RegistryValue", "RegistryValues"},
	}

	if _, ok := endings[currentSourceType]; !ok {
		return
	}

	trimmed := strings.TrimSpace(artifactDefinition.Name)
	if !strings.HasSuffix(trimmed, endings[currentSourceType][0]) && !strings.HasSuffix(trimmed, endings[currentSourceType][1]) {
		r.addCommon(filename, artifactDefinition.Name, "Artifact name should end in %s", strings.Join(endings[currentSourceType], " or "))
	}

}

func (r *validator) validateDocLong(filename string, artifactDefinition ArtifactDefinition) {
	if strings.Contains(artifactDefinition.Doc, "\n") && !strings.Contains(artifactDefinition.Doc, "\n\n") {
		r.addInfo(filename, artifactDefinition.Name, "Long docs should contain an empty line")
	}
}

func (r *validator) validateArtifactLabels(filename string, artifactDefinition ArtifactDefinition) {
	for _, labels := range artifactDefinition.Labels {
		found := false
		for _, l := range listLabels() {
			if l == labels {
				found = true
			}
		}
		if !found {
			r.addWarning(filename, artifactDefinition.Name, "Label %s is not valid", labels)
		}
	}
}

func (r *validator) validateArtifactOS(filename string, artifactDefinition ArtifactDefinition) {
	for _, supportedos := range artifactDefinition.SupportedOs {
		found := false
		for _, os := range listOSS() {
			if os == supportedos {
				found = true
			}
		}
		if !found {
			r.addWarning(filename, artifactDefinition.Name, "OS %s is not valid", supportedos)
		}
	}
}

func (r *validator) validateProvides(filename string, artifactDefinition ArtifactDefinition) {
	for _, provides := range artifactDefinition.Provides {
		if _, ok := knowledgeBase[provides]; !ok {
			r.addWarning(filename, artifactDefinition.Name, "Unused provides %s", provides)
		}
	}
}

func (r *validator) validateMacOSDoublePath(filename string, artifactDefinition ArtifactDefinition) {
	knownPaths := map[string]bool{}
	prefixes := []string{"/var", "/tmp", "/etc"}

	if isOSArtifactDefinition("Darwin", artifactDefinition.SupportedOs) {
		for _, source := range artifactDefinition.Sources {
			if isOSArtifactDefinition("Darwin", source.SupportedOs) {
				for _, path := range source.Attributes.Paths {
					for _, prefix := range prefixes {
						if strings.HasPrefix(path, prefix) || strings.HasPrefix(path, "/private"+prefix) {
							knownPaths[path] = true
						}
					}
				}
			}
		}
	}

	for knownPath := range knownPaths {
		var sibling string
		if strings.HasPrefix(knownPath, "/private") {
			sibling = strings.Replace(knownPath, "/private", "", 1)
		} else {
			sibling = "/private" + knownPath
		}
		if _, ok := knownPaths[sibling]; !ok {
			r.addWarning(filename, artifactDefinition.Name, "Found %s but not %s", knownPath, sibling)
		}
	}
}

// source

func (r *validator) validateUnnessesarryAttributes(filename, artifactDefinition string, source Source) {
	hasNames := len(source.Attributes.Names) > 0
	hasCommand := source.Attributes.Cmd != "" || len(source.Attributes.Args) > 0
	hasPaths := len(source.Attributes.Paths) > 0 || source.Attributes.Separator != ""
	hasKeys := len(source.Attributes.Keys) > 0
	hasKeyValuePairs := len(source.Attributes.KeyValuePairs) > 0
	hasWMI := source.Attributes.Query != "" || source.Attributes.BaseObject != ""

	switch source.Type {
	case sourceType.ArtifactGroup:
		if hasPaths || hasCommand || hasKeys || hasWMI || hasKeyValuePairs {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	case sourceType.Command:
		if hasNames || hasPaths || hasKeys || hasWMI || hasKeyValuePairs {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	case sourceType.Directory:
		fallthrough
	case sourceType.File:
		fallthrough
	case sourceType.Path:
		if hasNames || hasCommand || hasKeys || hasWMI || hasKeyValuePairs {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	case sourceType.RegistryKey:
		if hasNames || hasPaths || hasCommand || hasWMI || hasKeyValuePairs {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	case sourceType.RegistryValue:
		if hasNames || hasPaths || hasCommand || hasKeys || hasWMI {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	case sourceType.Wmi:
		if hasNames || hasPaths || hasCommand || hasKeys || hasKeyValuePairs {
			r.addWarning(filename, artifactDefinition, "Unnessesarry attribute set")
		}
	}
}
func (r *validator) validateRequiredAttributes(filename, artifactDefinition string, source Source) {
	switch source.Type {
	case sourceType.ArtifactGroup:
		if len(source.Attributes.Names) == 0 {
			r.addWarning(filename, artifactDefinition, "An ARTIFACT_GROUP requires the names attribute")
		}
	case sourceType.Command:
		if source.Attributes.Cmd == "" {
			r.addWarning(filename, artifactDefinition, "A COMMAND requires the cmd attribute")
		}
	}
}

func (r *validator) validateRequiredWindowsAttributes(filename, artifactDefinition string, source Source) {
	switch source.Type {
	case sourceType.Directory:
		fallthrough
	case sourceType.File:
		fallthrough
	case sourceType.Path:
		if len(source.Attributes.Paths) == 0 {
			r.addWarning(filename, artifactDefinition, "A %s requires the paths attribute", source.Type)
		}
		if source.Attributes.Separator != "" && source.Attributes.Separator != "\\" {
			r.addWarning(filename, artifactDefinition, "A %s requires a separator value of \"\\\" or \"\"", source.Type)
		}
	case sourceType.RegistryKey:
		if len(source.Attributes.Keys) == 0 {
			r.addWarning(filename, artifactDefinition, "A %s requires the keys attribute", source.Type)
		}
	case sourceType.RegistryValue:
		if len(source.Attributes.KeyValuePairs) == 0 {
			r.addWarning(filename, artifactDefinition, "A %s requires the key_value_pairs attribute", source.Type)
		}
	case sourceType.Wmi:
		if len(source.Attributes.Query) == 0 {
			r.addWarning(filename, artifactDefinition, "A %s requires the query attribute", source.Type)
		}
	}
}

func (r *validator) validateRequiredNonWindowsAttributes(filename, artifactDefinition string, source Source) {
	switch source.Type {
	case sourceType.Directory:
		fallthrough
	case sourceType.File:
		fallthrough
	case sourceType.Path:
		if len(source.Attributes.Paths) == 0 {
			r.addWarning(filename, artifactDefinition, "A %s requires the paths attribute", source.Type)
		}
	case sourceType.RegistryKey:
		fallthrough
	case sourceType.RegistryValue:
		fallthrough
	case sourceType.Wmi:
		r.addError(filename, artifactDefinition, "%s only supported for windows", source.Type)
	}
}

func (r *validator) validateRegistryCurrentControlSet(filename, artifactDefinition string, source Source) {

	err := `Registry key should not start with %%CURRENT_CONTROL_SET%%. Replace %%CURRENT_CONTROL_SET%% with HKEY_LOCAL_MACHINE\\System\\CurrentControlSet`

	for _, key := range source.Attributes.Keys {
		if strings.Contains(key, `%%CURRENT_CONTROL_SET%%`) {
			r.addInfo(filename, artifactDefinition, err)
		}
	}
	for _, keyvalue := range source.Attributes.KeyValuePairs {
		if strings.Contains(keyvalue.Key, `%%CURRENT_CONTROL_SET%%`) {
			r.addInfo(filename, artifactDefinition, err)
		}
	}
}

func (r *validator) validateRegistryHKEYCurrentUser(filename, artifactDefinition string, source Source) {
	err := `HKEY_CURRENT_USER\\ is not supported instead use: HKEY_USERS\\%%users.sid%%\\`
	for _, key := range source.Attributes.Keys {
		if strings.HasPrefix(key, `HKEY_CURRENT_USER\\`) {
			r.addError(filename, artifactDefinition, err)
		}
	}
	for _, keyvalue := range source.Attributes.KeyValuePairs {
		if strings.HasPrefix(keyvalue.Key, `HKEY_CURRENT_USER\\`) {
			r.addError(filename, artifactDefinition, err)
		}
	}
}

func (r *validator) validateDeprecatedVars(filename, artifactDefinition string, source Source) {
	deprecations := []struct {
		old, new string
	}{
		{old: "%%users.userprofile%%\\AppData\\Local", new: "%%users.localappdata%%"},
		{old: "%%users.userprofile%%\\AppData\\Roaming", new: "%%users.appdata%%"},
		{old: "%%users.userprofile%%\\Application Data", new: "%%users.appdata%%"},
		{old: "%%users.userprofile%%\\Local Settings\\Application Data", new: "%%users.localappdata%%"},
	}
	for _, path := range source.Attributes.Paths {
		for _, deprecation := range deprecations {
			if strings.Contains(path, deprecation.old) {
				r.addInfo(filename, artifactDefinition, "Replace %s by %s", deprecation.old, deprecation.new)
			}
		}
	}

}

// unc (r *validator) validateDoubleStar(filename, artifactDefinition string, source Source) {
// 	for _, path := range source.Attributes.Paths {
// 		if strings.Contains(path, `**`) {
// 			r.addInfo(filename, artifactDefinition, "Paths contains **")
// 			return
// 		}
// 	}
//

func (r *validator) validateNoWindowsHomedir(filename, artifactDefinition string, source Source) {
	windowsSource := len(source.SupportedOs) == 1 && source.SupportedOs[0] == supportedOS.Windows
	if len(source.SupportedOs) == 0 || windowsSource {
		for _, path := range source.Attributes.Paths {
			if strings.Contains(path, "%%users.homedir%%") {
				r.addInfo(filename, artifactDefinition, "Replace %s by %s", "%%users.homedir%%", "%%users.userprofile%%")
			}
		}
	}
}

func (r *validator) validateSourceType(filename, artifactDefinition string, source Source) {
	for _, t := range listTypes() {
		if t == source.Type {
			return
		}
	}
	r.addError(filename, artifactDefinition, "Type %s is not valid", source.Type)
}

func (r *validator) validateSourceOS(filename, artifactDefinition string, source Source) {
	for _, supportedos := range source.SupportedOs {
		found := false
		for _, os := range listOSS() {
			if os == supportedos {
				found = true
			}
		}
		if !found {
			r.addWarning(filename, artifactDefinition, "OS %s is not valid", supportedos)
		}
	}
}

func (r *validator) validateParameter(filename, artifactDefinition string, source Source) {

	FindInterpol := func(in string) (string, bool) {
		re := regexp.MustCompile(`%%.*?%%`)
		for _, match := range re.FindAllString(in, -1) {
			match = strings.Trim(match, `%`)
			if _, ok := knowledgeBase[match]; !ok {
				return match, false
			}
		}
		return "", true
	}

	for _, key := range source.Attributes.Keys {
		if match, found := FindInterpol(key); !found {
			r.addWarning(filename, artifactDefinition, "Parameter %s not found", match)
		}
	}
	for _, keyvalue := range source.Attributes.KeyValuePairs {
		if match, found := FindInterpol(keyvalue.Key); !found {
			r.addWarning(filename, artifactDefinition, "Parameter %s not found", match)

		}
	}
	for _, path := range source.Attributes.Paths {
		if match, found := FindInterpol(path); !found {
			r.addWarning(filename, artifactDefinition, "Parameter %s not found", match)

		}
	}

	if match, found := FindInterpol(source.Attributes.Query); !found {
		r.addWarning(filename, artifactDefinition, "Parameter %s not found", match)
	}
}
