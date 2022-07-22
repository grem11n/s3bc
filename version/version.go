package version

import "fmt"

// Set by ldflags
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
	BuiltBy = ""
)

// PrintInfo is used to get generic info about the build
func PrintInfo() {
	var result = fmt.Sprintf("s3bc:\nVersion: %s", Version)
	if Commit != "" {
		result = fmt.Sprintf("\n%s\nCommit: %s", result, Commit)
	}
	if Date != "" {
		result = fmt.Sprintf("\n%s\nBuilt at: %s", result, Date)
	}
	fmt.Println(result)
}
