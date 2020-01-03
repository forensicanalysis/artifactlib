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

package artifactlib_test

import (
	"fmt"

	"github.com/forensicanalysis/artifactlib/goartifacts"
	"github.com/forensicanalysis/fslib/filesystem/testfs"
)

func ExampleProcessArtifacts() {
	// This parses the arifact files, filters for the current OS, expands variables
	// and globing parameters

	// Create testing file system
	fs := &testfs.FS{}
	fs.CreateFile("foo.txt", nil)

	// []string{"Test1"}: Artifacts to filter for, artifacts groups are expanded
	// fs: File system used for expansion
	// true: Flag if multiple partitions are tried on windows
	// []string{"test/artifacts/collect_1.yaml"}: Files with artifact defintions
	artifacts, _ := goartifacts.ProcessFiles([]string{"Test1"}, fs, true, []string{"test/artifacts/collect_1.yaml"})

	// print resolved paths of the parsed artifact definition
	fmt.Println(artifacts[0].Sources[0].Attributes.Paths)
	// Output: [/foo.txt]
}

func ExampleValidate() {
	// Validate an artifact definition files
	flaws, _ := goartifacts.ValidateFiles([]string{"test/artifacts/collect_1.yaml"})

	// Print first issue
	fmt.Println(flaws[0].Message)
	// Output: Artifact name should end in File or Files
}
