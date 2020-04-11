// Copyright (c) 2019-2020 Siemens AG
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

// Package artifactlib provides a Go package and a Python library for processing
// forensic artifact definition files.
//
// Artifact definition files
//
// The artifact definition format is described in detail in the Style Guide (https://github.com/forensicanalysis/artifacts/blob/master/style_guide.md).
// The following shows an example for an artifact definition file. It defines the
// location of linux audit log files on a system.
//
// 	name: LinuxAuditLogFiles
// 	doc: Linux audit log files.
// 	sources:
// 	- type: FILE
// 	  attributes: {paths: ['/var/log/audit/*']}
// 	supported_os: [Linux]
//
// We use https://github.com/forensicanalysis/artifacts as the main repository for
// forensic artifacts definitions.
package artifactlib
