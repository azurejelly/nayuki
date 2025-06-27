package utils

import (
	"os"
	"runtime/debug"
)

const (
	GITHUB_REPOSITORY = "https://github.com/azurejelly/nayuki"
)

func ReadGitRevision() string {
	hash := ""

	if info, ok := debug.ReadBuildInfo(); ok {
		var dirty bool

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				hash = setting.Value
				if len(hash) > 8 {
					hash = hash[:8]
				}
			}

			if setting.Key == "vcs.modified" && setting.Value == "true" {
				dirty = true
			}
		}

		if dirty {
			hash += "-dirty"
		}
	}

	if hash == "" {
		return "unknown"
	}

	return hash
}

func IsDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}

	return true
}
