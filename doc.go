// Package artifactlib provides a Go package and a Python library for processing
// forensic artifact definition files.
//
// Artifact definition files
//
// The following shows an example for an artifact definition file. It defines the
// location of linux audit log files on a system.
//
// 	name: LinuxAuditLogFiles
// 	doc: Linux audit log files.
// 	sources:
// 	- type: FILE
// 	  attributes: {paths: ['/var/log/audit/*']}
// 	supported_os: [Linux]
//
// We use https://github.com/forensicanalysis/artifacts as the main repository for
// forensic artifacts definitions.
package artifactlib
