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
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// DecodeFile takes a single artifact definition file to decode.
func DecodeFile(filename string) ([]ArtifactDefinition, []Flaw, error) {
	var artifactDefinitions []ArtifactDefinition
	var flaws []Flaw

	// open file
	f, err := os.Open(filename)
	if err != nil {
		return artifactDefinitions, flaws, err
	}
	defer f.Close()

	// decode file
	dec := NewDecoder(f)
	artifactDefinitions, err = dec.Decode()
	if err != nil {
		if typeerror, ok := err.(*yaml.TypeError); ok {
			// parsing error
			for _, typeerr := range typeerror.Errors {
				flaws = append(flaws, Flaw{Error, typeerr, "", filename})
			}
		} else {
			// bad error
			return artifactDefinitions, flaws, err
		}
	}

	return artifactDefinitions, flaws, nil
}

// A Decoder reads and decodes YAML values from an input stream.
type Decoder struct {
	yamldecoder *yaml.Decoder
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read
// data from r beyond the YAML values requested.
func NewDecoder(r io.Reader) *Decoder {
	d := yaml.NewDecoder(r)
	d.SetStrict(true)
	return &Decoder{yamldecoder: d}
}

// Decode reads the next YAML-encoded value from its input and stores it in the
// value pointed to by v.
// See the documentation for Unmarshal for details about the conversion of YAML
// into a Go value.
func (dec *Decoder) Decode() ([]ArtifactDefinition, error) {
	var artifactDefinitions []ArtifactDefinition
	artifactDefinition := ArtifactDefinition{}
	for {
		// load every document
		err := dec.yamldecoder.Decode(&artifactDefinition)
		if err != nil {
			if err == io.EOF {
				return artifactDefinitions, nil
			}
			return artifactDefinitions, err
		}

		// gather artifact
		artifactDefinitions = append(artifactDefinitions, artifactDefinition)
	}
}
