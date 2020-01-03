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

// parameters lists all existing variables which can exist in some attributes
var knowledgeBase = map[string]string{
	"users.username":   "", // The name of the user."
	"users.temp":       "", // Temporary directory for the user."
	"users.desktop":    "", // The desktop directory of the user. E.g.  c:\\Documents and Settings\\foo\\Desktop",
	"users.last_logon": "", // RDFDatetime" description: "The last logon time for this user."
	"users.full_name":  "", // Full name of the user.

	// Windows specific values.
	"users.userdomain":       "", // The domain name of the user. E.g. MICROSOFT."
	"users.sid":              "", // The SID of the user as reported by the system. E.g.  S-1-5-80-859482183-879914841-863379149-1145462774-2388618682",
	"users.userprofile":      "", // The profile directory of the user. E.g. c:\\Documents and Settings\\foo.",
	"users.appdata":          "", // The %APPDATA% directory of the user as reported by the  system. E.g. c:\\Documents and Settings\\foo\\AppData\\Roaming",
	"users.localappdata":     "", // The %LOCALAPPDATA% directory of the user as reported by the " system. E.g. c:\\Documents and Settings\\foo\\AppData\\Local",
	"users.internet_cache":   "", // The cache directory of the user. E.g.  c:\\Documents and Settings\\foo\\AppData\\Local\\Temporary Internet Files.",
	"users.cookies":          "", // The cookies directory of the user. E.g.  c:\\Documents and Settings\\foo\\Cookies",
	"users.recent":           "", // The recent directory of the user. E.g.  c:\\Documents and Settings\\foo\\Recent.",
	"users.personal":         "", // The Personal directory of the user. E.g.  c:\\Documents and Settings\\foo\\Documents.",
	"users.startup":          "", // The Startup directory of the user. E.g.  c:\\Documents and Settings\\foo\\Startup.",
	"users.localappdata_low": "", // The LocalLow application data directory for data that " doesn't roam with the user.  E.g. %USERPROFILE%\\AppData\\LocalLow.  Vista and above.",

	// Posix specific values.
	"users.homedir":  "", // The homedir of the user as reported by the system. E.g.  "/home/foo",
	"users.uid":      "", // The uid of the of the user. E.g. 0."
	"users.gid":      "", // The gid of the the user. E.g. 5001."
	"users.shell":    "", // The shell of the the user. E.g. /bin/sh."
	"users.pw_entry": "", // The password state of the user, e.g.: shadow+sha512."
	// gids // Additional group ids the user is a member of."

	"fqdn":             "", // The fully qualified domain name reported by the OS. E.g.  host1.ad.foo.com",
	"time_zone":        "", // The timezone in Olson format E.g. Pacific/Galapagos.  http://en.wikipedia.org/wiki/Tz_database.",
	"os":               "", // The operating system. Case is important, must be one of  Windows Linux Darwin FreeBSD OpenBSD NetBSD",
	"os_major_version": "", // The major version of the OS, e.g. 7"
	"os_minor_version": "", // The minor version of the OS, e.g. 7"
	"environ_path":     "", // The system configured path variable."
	"environ_temp":     "", // The system temporary directory."

	// Linux specific distribution information.
	// See: lsb_release(1) man page, or the LSB Specification under the 'Command
	// Behaviour' section.
	"os_release": "", // Linux distribution name."

	// Windows specific system level parameters.
	"environ_systemroot":        "", // The value of the %SystemRoot% parameter, E.g. c:\\Windows"
	"environ_windir":            "", // The value of the %WINDIR% parameter. As returned by  the system, e.g. C:",
	"environ_programfiles":      "", // The value of the %PROGRAMFILES% parameter as returned by  the system, e.g. C:\\Program Files",
	"environ_programfilesx86":   "", // The value of the %PROGRAMFILES(X86)% parameter as returned " by the system, e.g. C:\\Program Files (x86)",
	"environ_systemdrive":       "", // The value of the %SystemDrive% parameter. As returned by  the system, e.g. C:",
	"environ_profilesdirectory": "", // Folder that typically contains users' profile directories;  e.g '%SystemDrive%\\Users'",
	"environ_allusersprofile":   "", // The value of the %AllUsersProfile% parameter. As returned  by the system, e.g. c:\\Documents and Settings\\All Users",
	"environ_allusersappdata":   "", // The value of the %AllUsersAppData% parameter. As returned  by the system, e.g. c:\\Documents and Settings\\All Users\\Application Data.",
	"current_control_set":       "", // The current value of the system CurrentControlSet  e.g. HKEY_LOCAL_MACHINE\\SYSTEM\\ControlSet001",
	"code_page":                 "", // The current code page of the system. Comes from  HKLM\\CurrentControlSet\\Control\\Nls\\CodePage e.g. cp1252.",
	"domain":                    "", // The domain the machine is connected to. E.g. MICROSOFT."
}

// sourceType is an enumeration of artifact definition source types
var sourceType = struct {
	ArtifactGroup string
	Command       string
	Directory     string
	File          string
	Path          string
	RegistryKey   string
	RegistryValue string
	Wmi           string
}{
	ArtifactGroup: "ARTIFACT_GROUP",
	Command:       "COMMAND",
	Directory:     "DIRECTORY",
	File:          "FILE",
	Path:          "PATH",
	RegistryKey:   "REGISTRY_KEY",
	RegistryValue: "REGISTRY_VALUE",
	Wmi:           "WMI",
}

// listTypes returns a list of all artifact definition source types
func listTypes() []string {
	return []string{
		sourceType.ArtifactGroup,
		sourceType.Command,
		sourceType.Directory,
		sourceType.File,
		sourceType.Path,
		sourceType.RegistryKey,
		sourceType.RegistryValue,
		sourceType.Wmi,
	}
}

// supportedOS is an enumeration of all supported OSs
var supportedOS = struct {
	Darwin  string
	Linux   string
	Windows string
}{
	Darwin:  "Darwin",
	Linux:   "Linux",
	Windows: "Windows",
}

func listOSS() []string {
	return []string{supportedOS.Darwin, supportedOS.Linux, supportedOS.Windows}
}

var label = struct {
	Antivirus          string
	Authentication     string
	Browser            string
	Cloud              string
	CloudStorage       string
	ConfigurationFiles string
	Docker             string
	ExternalMedia      string
	ExternalAccount    string
	Hadoop             string
	HistoryFiles       string
	Logs               string
	Mail               string
	Network            string
	Software           string
	System             string
	Users              string
	IOs                string
}{
	Antivirus:          "Antivirus",
	Authentication:     "Authentication",
	Browser:            "Browser",
	Cloud:              "Cloud",
	CloudStorage:       "Cloud Storage",
	ConfigurationFiles: "Configuration Files",
	Docker:             "Docker",
	ExternalMedia:      "External Media",
	ExternalAccount:    "ExternalAccount",
	Hadoop:             "Hadoop",
	HistoryFiles:       "History Files",
	Logs:               "Logs",
	Mail:               "Mail",
	Network:            "Network",
	Software:           "Software",
	System:             "System",
	Users:              "Users",
	IOs:                "iOS",
}

func listLabels() []string {
	return []string{
		label.Antivirus,
		label.Authentication,
		label.Browser,
		label.Cloud,
		label.CloudStorage,
		label.ConfigurationFiles,
		label.Docker,
		label.ExternalMedia,
		label.ExternalAccount,
		label.Hadoop,
		label.HistoryFiles,
		label.Logs,
		label.Mail,
		label.Network,
		label.Software,
		label.System,
		label.Users,
		label.IOs,
	}
}
