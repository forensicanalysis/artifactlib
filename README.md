<h1 align="center">artifactlib</h1>

<p  align="center">
 <a href="https://godocs.io/github.com/forensicanalysis/artifactlib/goartifacts"><img src="https://godocs.io/github.com/forensicanalysis/artifactlib?status.svg" alt="doc" /></a>
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
