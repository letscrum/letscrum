// package build contains the stuffs about build information such as: build time, git version and so on.
package build

import (
	"fmt"
	"runtime"
)

var (
	binVersion   string
	gitBranch    string
	gitTag       string
	gitCommit    string
	gitTreeState string
	buildDate    string
)

type version struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`

	GitBranch    string `json:"gitBranch"`
	GitTag       string `json:"gitTag"`
	GitTreeState string `json:"gitTreeState"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// Version returns the version information about the application.
func Version() *version {
	return &version{
		Version:   binVersion,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),

		GitBranch:    gitBranch,
		GitTag:       gitTag,
		GitTreeState: gitTreeState,
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
