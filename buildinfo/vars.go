package buildinfo

//go:generate /bin/sh buildinfo.sh

import (
	_ "embed"
	"strings"
)

//go:embed buildinfo.txt
var buildInfoRaw string

//nolint:gochecknoglobals // this is a static initialization.
var buildInfo struct {
	version, commitSHA, timestamp string
}

// Version for the app.
func Version() string {
	return buildInfo.version
}

// CommitSHA for the source.
func CommitSHA() string {
	return buildInfo.commitSHA
}

// Timestamp for the build.
func Timestamp() string {
	return buildInfo.commitSHA
}

//nolint:gochecknoinits // this is a static initialization.
func init() {
	values := strings.Split(buildInfoRaw, "\n")
	buildInfo.version = values[0]
	buildInfo.commitSHA = values[1]
	buildInfo.timestamp = values[2]
}
