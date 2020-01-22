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

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
)

func TestProcessFiles(t *testing.T) {
	result := map[string][]Source{
		"Test3Directory": {{
			Type: "DIRECTORY", Attributes: Attributes{Paths: []string{"/dev"}}, SupportedOs: []string{"Darwin", "Linux"},
		}},
	}

	if runtime.GOOS == "windows" {
		result = nil
	}

	type args struct {
		infs      fslib.FS
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]Source
		wantErr bool
	}{
		{"Valid Artifact Definitions", args{osfs.New(), []string{"../test/artifacts/valid/processing.yaml"}}, result, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &TestCollector{tt.args.infs, nil}

			err := ProcessFiles(nil, tt.args.filenames, collector)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(collector.Collected, tt.want) {
				t.Errorf("ProcessFiles() got = %#v, want %#v", collector.Collected, tt.want)
			}
			// if !reflect.DeepEqual(got1, tt.wantFlaws) {
			// 	t.Errorf("ProcessFiles() got1 = %v, want %v", got1, tt.wantFlaws)
			// }
		})
	}
}
