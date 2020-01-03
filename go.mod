module github.com/forensicanalysis/artifactlib

go 1.13

require (
	github.com/forensicanalysis/fslib v0.0.0-00010101000000-000000000000
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/looplab/tarjan v0.0.0-20161115091335-9cc6d6cebfb5
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.0-20180912035003-be2c049b30cc
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/forensicanalysis/fslib => ../fslib
