package version

import (
	"fmt"
	"runtime"
	"time"
)

// These variables will be set by the linker during build time
var (
	// Version is the semantic version of the CLI
	Version = "dev"

	// BuildTime is the time when the binary was built
	BuildTime = "unknown"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// GitBranch is the git branch name
	GitBranch = "unknown"

	// GoVersion is the Go runtime version
	GoVersion = runtime.Version()

	// Platform is the target platform (OS/Architecture)
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// GetVersion returns the semantic version
func GetVersion() string {
	return Version
}

// GetBuildTime returns the build time as a time.Time
func GetBuildTime() time.Time {
	if BuildTime == "unknown" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, BuildTime)
	return t
}

// GetGitCommit returns the git commit hash
func GetGitCommit() string {
	return GitCommit
}

// GetGitBranch returns the git branch name
func GetGitBranch() string {
	return GitBranch
}

// GetGoVersion returns the Go runtime version
func GetGoVersion() string {
	return GoVersion
}

// GetPlatform returns the target platform
func GetPlatform() string {
	return Platform
}

// BuildInfo returns a formatted string with all build information
func BuildInfo() string {
	versionType := GetVersionType()
	versionTypeEmoji := "📦"
	switch versionType {
	case "development":
		versionTypeEmoji = "🔧"
	case "prerelease":
		versionTypeEmoji = "🚧"
	case "stable":
		versionTypeEmoji = "✅"
	}

	return fmt.Sprintf(`🔖 Kavach CLI Version Information

%s Version:     %s (%s)
🕒 Build Time:  %s
🔗 Git Commit:  %s
🌿 Git Branch:  %s
⚡ Go Version:  %s
💻 Platform:    %s

For more information, visit: https://github.com/Gkemhcs/kavach-cli`,
		versionTypeEmoji, Version, versionType, BuildTime, GitCommit, GitBranch, GoVersion, Platform)
}

// ShortVersion returns just the version number
func ShortVersion() string {
	return Version
}

// IsDev returns true if this is a development build
func IsDev() bool {
	return Version == "dev" || Version == "unknown"
}

// IsRelease returns true if this is a release build
func IsRelease() bool {
	return !IsDev()
}

// IsPrerelease returns true if this is a pre-release (alpha, beta, rc)
func IsPrerelease() bool {
	return IsRelease() && (contains(Version, "alpha") || contains(Version, "beta") || contains(Version, "rc"))
}

// IsStable returns true if this is a stable release
func IsStable() bool {
	return IsRelease() && !IsPrerelease()
}

// GetVersionType returns the type of version
func GetVersionType() string {
	if IsDev() {
		return "development"
	}
	if IsPrerelease() {
		return "prerelease"
	}
	return "stable"
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
