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
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
)

func TestValidator_ValidateFiles(t *testing.T) {
	type args struct {
		yamlfile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Non existing file", args{"unknown.yaml"}, true},
		{"Valid Artifact Definitions", args{"../test/artifacts/valid/valid.yaml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidator()
			_, _, err := decodeFile(filepath.FromSlash(tt.args.yamlfile))
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(r.flaws) > 0 {
				t.Errorf("Validator.ValidateFiles() has flaws %v", r.flaws)
			}
		})
	}
}

func TestValidator_ValidateSyntax(t *testing.T) {
	type args struct {
		yamlfile string
	}
	tests := []struct {
		skipOnWindows bool
		name          string
		args          args
		want          []Flaw
	}{
		{true, "Non existing file", args{"unknown.yaml"}, []Flaw{{Error, "Error open unknown.yaml: no such file or directory", "", filepath.FromSlash("unknown.yaml")}}},
		{false, "Comment", args{"../test/artifacts/invalid/file_3.yaml"}, []Flaw{{Info, "The first line should be a comment", "", filepath.FromSlash("../test/artifacts/invalid/file_3.yaml")}}},
		{false, "Wrong file ending", args{"../test/artifacts/invalid/ending.yml"}, []Flaw{{Info, "File should have .yaml ending", "", filepath.FromSlash("../test/artifacts/invalid/ending.yml")}}},
		{false, "Whitespace at line end", args{"../test/artifacts/invalid/file_1.yaml"}, []Flaw{{Info, "Line 3 ends with whitespace", "", filepath.FromSlash("../test/artifacts/invalid/file_1.yaml")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnWindows && (runtime.GOOS == "windows") {
				t.Skip("OS and language specific error message is skipped")
			}
			r := newValidator()
			r.validateSyntax(filepath.FromSlash(tt.args.yamlfile))
			if !flawsEqual(r.flaws, tt.want) {
				t.Errorf("Validator.validateSyntax() = %#v, want %#v", r.flaws, tt.want)
			}
		})
	}
}

func TestValidator_ValidateFilesInvalid(t *testing.T) {
	type test struct {
		name     string
		yamlfile string
	}

	files, err := ioutil.ReadDir(filepath.Join("..", "test", "artifacts", "invalid"))
	if err != nil {
		t.Error(err.Error())
	}
	var tests []test
	for _, file := range files {
		tests = append(tests, test{"Test", filepath.Join("..", "test", "artifacts", "invalid", file.Name())})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newValidator()
			ads, flaws, err := decodeFile(tt.yamlfile)
			if err != nil {
				t.Error(err)
			}

			artifactDefinitionMap := map[string][]ArtifactDefinition{
				tt.yamlfile: ads,
			}

			r.validateArtifactDefinitions(artifactDefinitionMap)

			flaws = append(flaws, r.flaws...)
			if len(flaws) == 0 {
				t.Errorf("Validator.ValidateFiles() %s has no flaws", tt.yamlfile)
			}
		})
	}
}
