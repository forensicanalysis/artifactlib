<h1 align="center">artifactlib</h1>

<p  align="center">
 <a href="https://codecov.io/gh/forensicanalysis/artifactlib"><img src="https://codecov.io/gh/forensicanalysis/artifactlib/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://godocs.io/github.com/forensicanalysis/artifactlib"><img src="https://godocs.io/github.com/forensicanalysis/zipfs?status.svg" alt="doc" /></a>
</p>


The artifactlib project provides a Go package for processing
forensic artifact definition files.

### Example

```go
func main() {
	// Validate an artifact definition files
	flaws, _ := goartifacts.ValidateFiles([]string{"test/artifacts/collect_1.yaml"})

	// Print first issue
	fmt.Println(flaws[0].Message)
}
```
