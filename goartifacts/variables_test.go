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

/*
Test fails in gitlab CI for unix and will be replaced with new artifact definition

func Test_expandVar(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Errorf(err.Error())
	}
	userdir, err := homedir.Dir()
	if err != nil {
		t.Errorf(err.Error())
	}
	parts := strings.Split(currentUser.Username, "\\")

	type args struct {
		in string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		windowsOnly bool
	}{
		// {`Expand %%users.username%%`, args{`%%users.username%%`}, parts[len(parts)-1], true},
		// {`Expand %%users.sid%%`, args{`%%users.sid%%`}, currentUser.Uid, true},
		{`Expand %%users.homedir%%`, args{`%%users.homedir%%`}, userdir, false},
		{`Expand %%users.temp%%`, args{`%%users.temp%%`}, os.TempDir(), false},
		// {`Expand %%environ_windir%%`, args{`%%environ_windir%%`}, `C:\windows`, runtime.GOOS != "windows"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.windowsOnly || runtime.GOOS == "windows" {
				got := expandVar(tt.args.in)
				if !contains(got, tt.want) {
					t.Errorf("expandVar(%s) = %v, want %v", tt.args.in, got, tt.want)
				}
			}
		})
	}
}
*/
