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
	"reflect"
	"runtime"
	"testing"
)

func Test_filterOS(t *testing.T) {
	type args struct {
		artifactDefinitions []ArtifactDefinition
	}
	tests := []struct {
		name string
		args args
		want []ArtifactDefinition
	}{
		{"filterOS true", args{[]ArtifactDefinition{{Name: "Test", SupportedOs: []string{runtime.GOOS}}}}, []ArtifactDefinition{{Name: "Test", Sources: nil, SupportedOs: []string{runtime.GOOS}}}},
		{"filterOS sources", args{[]ArtifactDefinition{{Name: "Test", Sources: []Source{{SupportedOs: []string{runtime.GOOS}}, {SupportedOs: []string{"xxx"}}}}}}, []ArtifactDefinition{{Name: "Test", Sources: []Source{{SupportedOs: []string{runtime.GOOS}}}}}},
		{"filterOS false", args{[]ArtifactDefinition{{Name: "Test", SupportedOs: []string{"xxx"}}}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOS(tt.args.artifactDefinitions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterOS() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
