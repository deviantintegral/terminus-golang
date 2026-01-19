// Package version provides build-time version information for terminus-golang.
//
// Version information is injected at build time via ldflags. The following
// variables can be set:
//
//	-X github.com/deviantintegral/terminus-golang/pkg/version.version=1.0.0
//	-X github.com/deviantintegral/terminus-golang/pkg/version.commit=abc1234
//	-X github.com/deviantintegral/terminus-golang/pkg/version.dirty=true
//
// For tagged releases, set version to the tag (e.g., "0.6.0").
// For development builds, leave version empty and set commit to the short SHA.
// If the working tree is dirty, set dirty to "true".
package version

// These variables are set at build time via ldflags
var (
	// version is the semantic version (e.g., "0.6.0") for tagged releases.
	// Left empty for development builds.
	version string

	// commit is the short git commit hash (e.g., "abc1234").
	commit string

	// dirty indicates if the working tree had uncommitted changes.
	// Set to "true" if dirty, empty otherwise.
	dirty string
)

// String returns the version string for use in user agents and display.
//
// Returns:
//   - The semantic version if this is a tagged release (e.g., "0.6.0")
//   - The short commit hash if building from source (e.g., "abc1234")
//   - The commit hash with "-dirty" suffix if the tree was dirty (e.g., "abc1234-dirty")
//   - "dev" if no version information is available
func String() string {
	// If we have a version tag, use it (tagged release)
	if version != "" {
		return version
	}

	// If we have a commit hash, use it (development build)
	if commit != "" {
		if dirty == "true" {
			return commit + "-dirty"
		}
		return commit
	}

	// Fallback for local development without ldflags
	return "dev"
}

// Commit returns the short git commit hash, or empty string if not set.
func Commit() string {
	return commit
}

// IsDirty returns true if the build was from a dirty working tree.
func IsDirty() bool {
	return dirty == "true"
}

// IsRelease returns true if this is a tagged release build.
func IsRelease() bool {
	return version != ""
}
