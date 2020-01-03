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

// +build !windows

package goartifacts

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

func getUsers() (users []user.User, err error) {
	var Users []string
	file, err := os.Open("/etc/passwd")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		// skip all line starting with #
		if equal := strings.Index(line, "#"); equal < 0 {
			// get the username and description
			lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
				return divide == ':' // we divide at colon
			})

			if len(lineSlice) > 0 {
				Users = append(Users, lineSlice[0])
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	for _, name := range Users {
		usr, err := user.Lookup(name)
		if err != nil {
			return nil, err
		}
		users = append(users, *usr)
	}
	return
}

func getUsernames(string) (usernames []string, err error) {
	users, err := getUsers()
	for _, user := range users {
		usernames = append(usernames, user.Name)
	}
	return
}

func getHomedirs(string) (homedirs []string, err error) {
	users, err := getUsers()
	for _, user := range users {
		homedirs = append(homedirs, user.HomeDir)
	}
	return
}

func getSIDs(string) (uids []string, err error) {
	users, err := getUsers()
	for _, user := range users {
		uids = append(uids, user.Gid)
	}
	return
}
