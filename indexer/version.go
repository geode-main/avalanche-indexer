package indexer

import "fmt"

var (
	AppName    = "avalanche-indexer"
	AppVersion = "0.6.0"
	GitCommit  = "-"
	GoVersion  = "-"
)

func VersionString() string {
	return fmt.Sprintf(
		"%s %s git=%q go=%q",
		AppName,
		AppVersion,
		GitCommit,
		GoVersion,
	)
}
