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
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Severity level of a flaw.
type Severity int

// Severity levels of a flaw.
const (
	Common  Severity = iota // Common errors
	Info                    // Style violations, will not create any issues
	Warning                 // Will compile but might create unexpected results
	Error                   // Will likely become an error
)

// Flaw is a single issue found by the validator
type Flaw struct {
	Severity           Severity
	Message            string
	ArtifactDefinition string
	File               string
}

// The validator performs all validations and stores the found flaws
type validator struct {
	flaws []Flaw
}

func newValidator() *validator {
	return &validator{[]Flaw{}}
}

func (r *validator) addFlaw(filename, artifactDefiniton string, severity Severity, format string, a ...interface{}) {
	r.flaws = append(r.flaws, Flaw{severity, fmt.Sprintf(format, a...), artifactDefiniton, filename})
}
func (r *validator) addCommon(filename, artifactDefiniton, format string, a ...interface{}) {
	r.addFlaw(filename, artifactDefiniton, Common, format, a...)
}
func (r *validator) addInfo(filename, artifactDefiniton, format string, a ...interface{}) {
	r.addFlaw(filename, artifactDefiniton, Info, format, a...)
}
func (r *validator) addWarning(filename, artifactDefiniton, format string, a ...interface{}) {
	r.addFlaw(filename, artifactDefiniton, Warning, format, a...)
}
func (r *validator) addError(filename, artifactDefiniton, format string, a ...interface{}) {
	r.addFlaw(filename, artifactDefiniton, Error, format, a...)
}

// ValidateArtifactDefinitions validates a map of artifact definitions and returns any flaws found in those.
func ValidateArtifactDefinitions(artifactDefinitionMap map[string][]ArtifactDefinition) []Flaw {
	r := newValidator()
	r.validateArtifactDefinitions(artifactDefinitionMap)
	return r.flaws
}

// validate single artifacts
func (r *validator) validateArtifactDefinitions(artifactDefinitionMap map[string][]ArtifactDefinition) {
	var globalArtifactDefinitions []ArtifactDefinition

	for filename, artifactDefinitions := range artifactDefinitionMap {
		if filename != "" {
			r.validateSyntax(filename)
		}

		globalArtifactDefinitions = append(globalArtifactDefinitions, artifactDefinitions...)
		for _, artifactDefinition := range artifactDefinitions {
			r.validateArtifactDefinition(filename, artifactDefinition)
		}
	}

	// global validations
	r.validateNameUnique(globalArtifactDefinitions)
	r.validateRegistryKeyUnique(globalArtifactDefinitions)
	r.validateRegistryValueUnique(globalArtifactDefinitions)
	r.validateGroupMemberExist(globalArtifactDefinitions)
	r.validateNoCycles(globalArtifactDefinitions)
	r.validateParametersProvided(globalArtifactDefinitions)

}

// validate single artifacts
func (r *validator) validateArtifactDefinition(filename string, artifactDefinition ArtifactDefinition) {
	r.validateNameCase(filename, artifactDefinition)
	r.validateNameTypeSuffix(filename, artifactDefinition)
	r.validateDocLong(filename, artifactDefinition)
	r.validateNamePrefix(filename, artifactDefinition)
	r.validateOSSpecific(filename, artifactDefinition)
	r.validateArtifactOS(filename, artifactDefinition)
	r.validateArtifactLabels(filename, artifactDefinition)
	r.validateProvides(filename, artifactDefinition)
	if isOSArtifactDefinition(supportedOS.Darwin, artifactDefinition.SupportedOs) {
		r.validateMacOSDoublePath(filename, artifactDefinition)
	}

	// validate sources
	for _, source := range artifactDefinition.Sources {
		r.validateUnnessesarryAttributes(filename, artifactDefinition.Name, source)
		r.validateRequiredAttributes(filename, artifactDefinition.Name, source)
		r.validateDeprecatedVars(filename, artifactDefinition.Name, source)
		r.validateRegistryCurrentControlSet(filename, artifactDefinition.Name, source)
		r.validateRegistryHKEYCurrentUser(filename, artifactDefinition.Name, source)
		// r.validateDoubleStar(filename, artifactDefinition.Name, source)
		r.validateSourceOS(filename, artifactDefinition.Name, source)
		r.validateSourceType(filename, artifactDefinition.Name, source)
		r.validateParameter(filename, artifactDefinition.Name, source)

		if isOSArtifactDefinition(supportedOS.Windows, artifactDefinition.SupportedOs) && isOSArtifactDefinition(supportedOS.Windows, source.SupportedOs) {
			r.validateNoWindowsHomedir(filename, artifactDefinition.Name, source)
			r.validateRequiredWindowsAttributes(filename, artifactDefinition.Name, source)
		}
		if (isOSArtifactDefinition(supportedOS.Linux, artifactDefinition.SupportedOs) || isOSArtifactDefinition(supportedOS.Darwin, artifactDefinition.SupportedOs)) &&
			(isOSArtifactDefinition(supportedOS.Linux, source.SupportedOs) || isOSArtifactDefinition(supportedOS.Darwin, source.SupportedOs)) {
			r.validateRequiredNonWindowsAttributes(filename, artifactDefinition.Name, source)
		}
	}
}

func (r *validator) validateSyntax(filename string) {
	if !strings.HasSuffix(filename, ".yaml") {
		r.addInfo(filename, "", "File should have .yaml ending")
	}

	// open file
	f, err := os.Open(filename)
	if err != nil {
		r.addError(filename, "", "Error %s", err)
		return
	}
	defer f.Close()
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 {
			if len(line) < 3 || !strings.HasPrefix(line, "# ") {
				r.addInfo(filename, "", "The first line should be a comment")
			}
		}

		if line != strings.TrimRight(line, " \t") {
			r.addInfo(filename, "", "Line %d ends with whitespace", i+1)
		}
		i++
	}
}
