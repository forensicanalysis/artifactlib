<h1 align="center">artifactlib</h1>

<p  align="center">
 <a href="https://github.com/forensicanalysis/artifactlib/actions"><img src="https://github.com/forensicanalysis/artifactlib/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/artifactlib"><img src="https://codecov.io/gh/forensicanalysis/artifactlib/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/artifactlib"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/artifactlib" alt="report" /></a>
 <a href="https://pkg.go.dev/github.com/forensicanalysis/artifactlib"><img src="https://img.shields.io/badge/go.dev-documentation-007d9c?logo=go&logoColor=white" alt="doc" /></a>
</p>


The artifactlib project provides a Go package and a Python library for processing
forensic artifact definition files.

## Artifact definition files
The artifact definition format is described in detail in the Style Guide ([https://github.com/forensicanalysis/artifacts/blob/master/style_guide.md](https://github.com/forensicanalysis/artifacts/blob/master/style_guide.md)).
The following shows an example for an artifact definition file. It defines the
location of linux audit log files on a system.

```
name: LinuxAuditLogFiles
doc: Linux audit log files.
sources:
- type: FILE
  attributes: {paths: ['/var/log/audit/*']}
supported_os: [Linux]
```

We use [https://github.com/forensicanalysis/artifacts](https://github.com/forensicanalysis/artifacts) as the main repository for
forensic artifacts definitions.


## Python library

### Installation

Python installation can be easily done via pip:

```bash
pip install pyartifacts
```

### Usage

```python
from pyartifacts.registry import Registry

if __name__ == '__main__':
    registry = Registry()
    registry.read_folder("test/artifacts/valid")
    print(registry)
```

The full documentation can be found in [the documentation](https://forensicanalysis.github.io/artifactlib/pyartifacts/docs/html).

## Go package

### Installation


```bash
go get -u github.com/forensicanalysis/artifactlib
```


### Usage


### Validate
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

## Acknowledgment

The development of this software was partially sponsored by Siemens CERT, but
is not an official Siemens product.
