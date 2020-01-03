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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func expandWinEnv(env string) ([]string, error) {
	env = strings.Trim(env, "%")
	env = strings.Replace(env, "environ_", "", 1)
	env = strings.ToUpper(env)
	val, ok := os.LookupEnv(env)
	if !ok {
		return []string{}, fmt.Errorf("environment variable %s could not be resolved", env)
	}
	return []string{val}, nil
}

func replaces(s, old string, news []string) (ss []string) {
	for _, newString := range news {
		ss = append(ss, strings.ReplaceAll(s, old, newString))
	}
	return
}

func expandTmp(_ string) ([]string, error) {
	return []string{os.TempDir()}, nil
}

func usersData(v string) ([]string, error) {
	switch v {
	case `%%users.appdata%%`:
		return []string{`%%users.homedir%%/AppData/Roaming`}, nil
	case `%%users.localappdata%%`:
		return []string{`%%users.homedir%%/AppData/Local`}, nil
	}
	return nil, errors.New("variable does not exist")
}

func expandVar(ins string) (out []string) {
	vars := []struct {
		placeholder string
		fun         func(string) ([]string, error)
	}{
		{placeholder: `%%users.appdata%%`, fun: usersData},
		{placeholder: `%%users.localappdata%%`, fun: usersData},

		{placeholder: `%%users.sid%%`, fun: getSIDs},
		{placeholder: `%%users.userprofile%%`, fun: getHomedirs},
		{placeholder: `%%users.homedir%%`, fun: getHomedirs},
		{placeholder: `%%users.username%%`, fun: getUsernames},

		{placeholder: `%%users.temp%%`, fun: expandTmp},

		{placeholder: `%%environ_systemroot%%`, fun: expandWinEnv},
		{placeholder: `%%environ_windir%%`, fun: expandWinEnv},
		{placeholder: `%%environ_programfiles%%`, fun: expandWinEnv},
		{placeholder: `%%environ_programfilesx86%%`, fun: expandWinEnv},
		{placeholder: `%%environ_systemdrive%%`, fun: expandWinEnv},
		{placeholder: `%%environ_allusersprofile%%`, fun: expandWinEnv},
		{placeholder: `%%environ_allusersappdata%%`, fun: expandWinEnv},

		// {placeholder: `%%users.desktop%%`, fun: expandError},
		// {placeholder: `%%users.last_logon%%`, fun: expandError},
		// {placeholder: `%%users.full_name%%`, fun: expandUserVar},
		// {placeholder: `%%users.userdomain%%`, fun: expandError},
		// {placeholder: `%%users.internet_cache%%`, fun: expandError},
		// {placeholder: `%%users.cookies%%`, fun: expandError},
		// {placeholder: `%%users.recent%%`, fun: expandError},
		// {placeholder: `%%users.personal%%`, fun: expandError},
		// {placeholder: `%%users.startup%%`, fun: expandError},
		// {placeholder: `%%users.localappdata_low%%`, fun: expandError},
		// {placeholder: `%%users.uid%%`, fun: expandUserVar},
		// {placeholder: `%%users.gid%%`, fun: expandUserVar},
		// {placeholder: `%%users.shell%%`, fun: expandError},
		// {placeholder: `%%users.pw_entry%%`, fun: expandError},
		// {placeholder: `%%fqdn%%`, fun: expandError},
		// {placeholder: `%%time_zone%%`, fun: expandError},
		// {placeholder: `%%os%%`, fun: expandOS},
		// {placeholder: `%%os_major_version%%`, fun: expandError},
		// {placeholder: `%%os_minor_version%%`, fun: expandError},
		// {placeholder: `%%environ_path%%`, fun: expandError},
		// {placeholder: `%%environ_temp%%`, fun: expandError},
		// {placeholder: `%%os_release%%`, fun: expandError},
		// {placeholder: `%%environ_profilesdirectory%%`, fun: expandWinEnv},
		// {placeholder: `%%current_control_set%%`, fun: expandError},
		// {placeholder: `%%code_page%%`, fun: expandError},
		// {placeholder: `%%domain%%`, fun: expandError},
	}
	in := []string{ins}
	for _, v := range vars {
		for _, i := range in {
			if strings.Contains(i, v.placeholder) {
				newvalues, err := v.fun(v.placeholder)
				if err != nil {
					log.Println(err)
				}
				out = append(out, replaces(i, v.placeholder, newvalues)...)
			} else {
				out = append(out, i)
			}
		}
		in = out
		out = nil
	}
	return in
}
