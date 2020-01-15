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

// Package artifactvalidator implements the artifactvalidator command line tool
// that can validate artifact definition files and search for errors, possible
// inconsistencies and other flaws.
//
// Usage
//
// To run just provide the location of the forensic artifact definition files:
//     artifactvalidator -v -s artifacts/data/*.yaml
// The output is a list of potential issues in those files.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	log "github.com/inconshreveable/log15"
	"github.com/olekukonko/tablewriter"

	"github.com/forensicanalysis/artifactlib/goartifacts"
)

var logger = log.New()

func main() {
	exitcode := 0
	// parse flags
	var verbose, summary, quite, nofail bool
	flag.BoolVar(&verbose, "verbose", false, "show common flaws as well")
	flag.BoolVar(&verbose, "v", false, "show common flaws as well"+" (shorthand)")
	flag.BoolVar(&quite, "quite", false, "hide informational flaws")
	flag.BoolVar(&quite, "q", false, "hide informational flaws"+" (shorthand)")
	flag.BoolVar(&summary, "summary", false, "show summary")
	flag.BoolVar(&summary, "s", false, "show summary"+" (shorthand)")
	flag.BoolVar(&nofail, "no-fail", false, "do not fail on flaws")
	flag.Parse()

	// setup logging
	handler := log.StreamHandler(os.Stdout, log.TerminalFormat())
	switch {
	case verbose:
		logger.SetHandler(handler)
	case quite:
		logger.SetHandler(log.LvlFilterHandler(log.LvlWarn, handler))
	default:
		logger.SetHandler(log.LvlFilterHandler(log.LvlInfo, handler))
	}

	args := flag.Args()

	// windows does not expand *
	if runtime.GOOS == "windows" {
		var files []string
		for _, arg := range args {
			paths, err := filepath.Glob(arg)
			if err != nil {
				logger.Error(err.Error())
			}
			files = append(files, paths...)
		}
		args = files
	}

	// parse artifacts
	flaws, err := goartifacts.ValidateFiles(args)
	if err != nil {
		logger.Crit(err.Error())
		os.Exit(1)
	}

	var filteredFlaws []goartifacts.Flaw
	if verbose {
		filteredFlaws = flaws
	} else {
		for _, flaw := range flaws {
			if flaw.Severity >= goartifacts.Warning || (!quite && flaw.Severity == goartifacts.Info) {
				filteredFlaws = append(filteredFlaws, flaw)
			}
		}
	}

	if len(filteredFlaws) > 0 {
		exitcode = 1
		printFlaws(filteredFlaws)
		if nofail {
			exitcode = 0
		}
	}

	if summary {
		// decode file
		var artifactDefinitions []goartifacts.ArtifactDefinition
		for _, filename := range args {
			ads, _, err := goartifacts.DecodeFile(filename)
			if err != nil {
				logger.Crit(err.Error())
				os.Exit(1)
			}
			artifactDefinitions = append(artifactDefinitions, ads...)
		}

		fmt.Printf("\nFound %d artifacts\n", len(artifactDefinitions))

		if len(artifactDefinitions) > 0 {
			var sourcetypes, oss, labels = map[string]int{}, map[string]int{}, map[string]int{}
			for _, artifactDefinition := range artifactDefinitions {
				for _, source := range artifactDefinition.Sources {
					inc(sourcetypes, source.Type)
				}
				for _, supportedOS := range artifactDefinition.SupportedOs {
					inc(oss, supportedOS)
				}
				for _, label := range artifactDefinition.Labels {
					inc(labels, label)
				}
			}
			printTable("Artifact definition by type", sourcetypes)
			printTable("Artifact definition by OS", oss)
			printTable("Artifact definition by label", labels)
		}
	}
	os.Exit(exitcode)
}

func inc(m map[string]int, key string) {
	if _, ok := m[key]; !ok {
		m[key] = 0
	}
	m[key]++
}

func printTable(caption string, m map[string]int) {
	fmt.Println("\n" + caption)
	keys, values := sortedMap(m)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(keys)
	table.SetCenterSeparator("|")
	table.AppendBulk([][]string{values})
	table.Render()
}

func sortedMap(m map[string]int) ([]string, []string) {
	var values []string
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		values = append(values, fmt.Sprint(m[k]))
	}
	return keys, values
}

func printFlaws(flaws []goartifacts.Flaw) {
	for _, flaw := range flaws {
		switch flaw.Severity {
		case goartifacts.Common:
			logger.Debug(fmt.Sprintf("%-60s %-30s %s", flaw.File, flaw.ArtifactDefinition, flaw.Message))
		case goartifacts.Info:
			logger.Info(fmt.Sprintf("%-60s %-30s %s", flaw.File, flaw.ArtifactDefinition, flaw.Message))
		case goartifacts.Warning:
			logger.Warn(fmt.Sprintf("%-60s %-30s %s", flaw.File, flaw.ArtifactDefinition, flaw.Message))
		case goartifacts.Error:
			logger.Error(fmt.Sprintf("%-60s %-30s %s", flaw.File, flaw.ArtifactDefinition, flaw.Message))
		}
	}
}
