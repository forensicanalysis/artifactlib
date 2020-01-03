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
	"github.com/forensicanalysis/fslib"
	"reflect"
	"runtime"
	"testing"

	"github.com/forensicanalysis/fslib/filesystem/osfs"
)

func TestProcessFiles(t *testing.T) {
	result := []ArtifactDefinition{{
		Name: "TestDirectory",
		Doc:  "Minimal dummy artifact definition for tests",
		Sources: []Source{{
			Type: "DIRECTORY", Attributes: Attributes{Paths: []string{"/etc"}}, SupportedOs: []string{"Darwin", "Linux"},
		}},
	}}

	if runtime.GOOS == "windows" {
		result[0].Sources = nil
	}

	type args struct {
		infs      fslib.FS
		filenames []string
	}
	tests := []struct {
		name      string
		args      args
		want      []ArtifactDefinition
		wantFlaws []Flaw
		wantErr   bool
	}{
		{"Valid Artifact Definitions", args{osfs.New(), []string{"../test/artifacts/valid/processing.yaml"}}, result, []Flaw{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessFiles(nil, tt.args.infs, false, tt.args.filenames)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessFiles() got = %#v, want %#v", got, tt.want)
			}
			// if !reflect.DeepEqual(got1, tt.wantFlaws) {
			// 	t.Errorf("ProcessFiles() got1 = %v, want %v", got1, tt.wantFlaws)
			// }
		})
	}
}

func TestDecodeFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []ArtifactDefinition
		want1   []Flaw
		wantErr bool
	}{
		{"Valid Artifact Definitions", args{"../test/artifacts/valid/mac_os_double_path_3.yaml"}, []ArtifactDefinition{{Name: "TestDirectory", Doc: "Minimal dummy artifact definition for tests", Sources: []Source{{Type: "DIRECTORY", Attributes: Attributes{Paths: []string{"/etc", "/private/etc"}}, SupportedOs: []string{"Darwin"}}}}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeFile() got = %#v, want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DecodeFile() got1 = %#v, want %#v", got1, tt.want1)
			}
		})
	}
}

func Test_isOSArtifactDefinition(t *testing.T) {
	type args struct {
		os          string
		supportedOs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test Windows", args{"Windows", []string{"Windows"}}, true},
		{"Test Windows", args{"Windows", []string{"Linux", "Darwin"}}, false},
		{"Test Linux", args{"Linux", []string{"Linux"}}, true},
		{"Test Darwin", args{"Darwin", []string{"Darwin"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOSArtifactDefinition(tt.args.os, tt.args.supportedOs); got != tt.want {
				t.Errorf("isOSArtifactDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}
