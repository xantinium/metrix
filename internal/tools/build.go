package tools

import "fmt"

// BuildInfo информация о сборке.
type BuildInfo struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

// PrintBuildInfo выводит информацию о сборке.
func PrintBuildInfo(info BuildInfo) {
	const notApplicable = "N/A"

	buildVersion := info.BuildVersion
	if buildVersion == "" {
		buildVersion = notApplicable
	}

	buildDate := info.BuildDate
	if buildDate == "" {
		buildDate = notApplicable
	}

	buildCommit := info.BuildCommit
	if buildCommit == "" {
		buildCommit = notApplicable
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
