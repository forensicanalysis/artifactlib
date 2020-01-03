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
	"sort"
	"testing"
)

func Test_expandKey(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Expand Star", args{"/*"}, []string{"/HKEY_CLASSES_ROOT", "/HKEY_CURRENT_USER", "/HKEY_LOCAL_MACHINE", "/HKEY_USERS", "/HKEY_CURRENT_CONFIG"}},
		{"Expand Key", args{"/NOKEY"}, []string{}},
		{"Expand HKEY_LOCAL_MACHINE star", args{`/HKEY_LOCAL_MACHINE/*`}, []string{"/HKEY_LOCAL_MACHINE/HARDWARE", "/HKEY_LOCAL_MACHINE/SAM", "/HKEY_LOCAL_MACHINE/SOFTWARE", "/HKEY_LOCAL_MACHINE/SYSTEM"}},
		{"Expand HKEY_LOCAL_MACHINE double star", args{`/HKEY_LOCAL_MACHINE/**`}, []string{"/HKEY_LOCAL_MACHINE/HARDWARE", "/HKEY_LOCAL_MACHINE/SYSTEM/CurrentControlSet/Control"}}, // any many many more keys
		{"Expand CurrentControlSet star", args{`/HKEY_LOCAL_MACHINE/System/CurrentControlSet/*`}, []string{"/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control", "/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Enum", "/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Hardware Profiles", "/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Policies", "/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Services", "/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Software"}},
		{"Expand ComputerName", args{`/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName`}, []string{`/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandKey(tt.args.s)
			if err != nil {
				t.Error(err)
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !isSubset(got, tt.want) {
				t.Errorf("expandKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
