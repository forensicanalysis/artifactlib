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
	"github.com/forensicanalysis/fslib/filesystem/registryfs"
	"strings"
)

func getUsernames(_ string) (usernames []string, err error) {
	regfs := registryfs.New()
	sids, err := getSIDs("")
	if err != nil {
		return nil, err
	}

	for _, sid := range sids {
		f, err := regfs.Open("/HKEY_LOCAL_MACHINE/SOFTWARE/Microsoft/Windows NT/CurrentVersion/ProfileList/" + sid)
		if err != nil {
			return nil, err
		}

		value, _, err := f.(*registryfs.Key).Key.GetStringValue("ProfileImagePath")
		if err != nil {
			return nil, err
		}
		parts := strings.Split(value, "\\")

		usernames = append(usernames, parts[len(parts)-1])
	}
	return
}

func getHomedirs(_ string) (homedirs []string, err error) {
	regfs := registryfs.New()
	sids, err := getSIDs("")
	if err != nil {
		return nil, err
	}

	for _, sid := range sids {
		f, err := regfs.Open("/HKEY_LOCAL_MACHINE/SOFTWARE/Microsoft/Windows NT/CurrentVersion/ProfileList/" + sid)
		if err != nil {
			return nil, err
		}

		value, _, err := f.(*registryfs.Key).Key.GetStringValue("ProfileImagePath")
		if err != nil {
			return nil, err
		}

		homedirs = append(homedirs, value)
	}
	return
}

func getSIDs(_ string) (sids []string, err error) {
	regfs := registryfs.New()
	f, err := regfs.Open("/HKEY_LOCAL_MACHINE/SOFTWARE/Microsoft/Windows NT/CurrentVersion/ProfileList")
	if err != nil {
		return nil, err
	}
	return f.Readdirnames(0)
}
