<h1 align="center">artifactlib</h1>

<p  align="center">
 <a href="https://github.com/forensicanalysis/artifactlib/actions"><img src="https://github.com/forensicanalysis/artifactlib/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/artifactlib"><img src="https://codecov.io/gh/forensicanalysis/artifactlib/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/artifactlib"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/artifactlib" alt="report" /></a>
 <a href="https://godocs.io/github.com/forensicanalysis/artifactlib"><img src="https://godocs.io/github.com/forensicanalysis/zipfs?status.svg" alt="doc" /></a>
 <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Fartifactlib?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Fartifactlib.svg?type=shield"/></a>
</p>


The artifactlib project provides a Go package for processing
forensic artifact definition files.

### Installation

```bash
go get -u github.com/forensicanalysis/artifactlib
```

### Example

```go
package main

import (
	"fmt"
	"github.com/forensicanalysis/artifactlib/goartifacts"
)

func main() {
	// Validate an artifact definition files
	flaws, _ := goartifacts.ValidateFiles([]string{"test/artifacts/collect_1.yaml"})

	// Print first issue
	fmt.Println(flaws[0].Message)
}

```

## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
