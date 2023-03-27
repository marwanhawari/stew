package stew

import (
	"fmt"

	"github.com/marwanhawari/stew/constants"
)

// ReleasesNotFoundError occurs if no releases are found for a GitHub repo.
type ReleasesNotFoundError struct {
	Owner string
	Repo  string
}

func (e ReleasesNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find any releases for %v", constants.RedColor("Error:"), constants.RedColor("https://github.com/"+e.Owner+"/"+e.Repo))
}

// AssetsNotFoundError occurs if no assets are found for a GitHub release.
type AssetsNotFoundError struct {
	Tag string
}

func (e AssetsNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find any assets for release %v", constants.RedColor("Error:"), constants.RedColor(e.Tag))
}

// NoPackagesInLockfileError occurs if you try to remove packages from a lockfile without any packages.
type NoPackagesInLockfileError struct{}

func (e NoPackagesInLockfileError) Error() string {
	return fmt.Sprintf("%v Cannot remove from an empty packages slice in the lockfile", constants.RedColor("Error:"))
}

// IndexOutOfBoundsInLockfileError occurs if you try to access an out-of-bounds index in the lockfile packages.
type IndexOutOfBoundsInLockfileError struct{}

func (e IndexOutOfBoundsInLockfileError) Error() string {
	return fmt.Sprintf("%v Index out of bounds in lockfile packages", constants.RedColor("Error:"))
}

// StewpathNotFoundError occurs if the ~/.stew path is not found.
type StewpathNotFoundError struct {
	StewPath string
}

func (e StewpathNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find the stew path at %v", constants.RedColor("Error:"), constants.RedColor(e.StewPath))
}

// NonZeroStatusCodeDownloadError occurs if a non-zero status code is received when trying to download a file.
type NonZeroStatusCodeDownloadError struct {
	StatusCode int
}

func (e NonZeroStatusCodeDownloadError) Error() string {
	return fmt.Sprintf("%v Received non-zero status code from HTTP request when attempting to download a file: %v", constants.RedColor("Error:"), constants.RedColor(e.StatusCode))
}

// EmptyCLIInputError occurs if the CLI input is empty.
type EmptyCLIInputError struct{}

func (e EmptyCLIInputError) Error() string {
	return fmt.Sprintf("%v Input cannot be empty. Use the --help flag for more info", constants.RedColor("Error:"))
}

// CLIFlagAndInputError occurs if you try to use a CLI flag with a CLI input at the same time.
type CLIFlagAndInputError struct{}

func (e CLIFlagAndInputError) Error() string {
	return fmt.Sprintf("%v Cannot use the --all flag with a positional argument", constants.RedColor("Error:"))
}

// AssetAlreadyDownloadedError occurs if the requested asset has already been downloaded.
type AssetAlreadyDownloadedError struct {
	Asset string
}

func (e AssetAlreadyDownloadedError) Error() string {
	return fmt.Sprintf("%v The %v asset has already been downloaded and installed", constants.RedColor("Error:"), constants.RedColor(e.Asset))
}

// CouldntDetectBinaryError occurs if the binary could not be automatically detected.
type CouldntDetectBinaryError struct {
	Options []string
}

func (e CouldntDetectBinaryError) Error() string {
	return fmt.Sprintf("%v Could not automatically detect the binary. Options: %v", constants.RedColor("Error:"), constants.RedColor(e.Options))
}

// AbortBinaryOverwriteError occurs if the overwrite of a binary is aborted.
type AbortBinaryOverwriteError struct {
	Binary string
}

func (e AbortBinaryOverwriteError) Error() string {
	return fmt.Sprintf("%v Overwrite of %v aborted", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

// BinaryNotInstalledError occurs if you try to operate on a binary that is not installed.
type BinaryNotInstalledError struct {
	Binary string
}

func (e BinaryNotInstalledError) Error() string {
	return fmt.Sprintf("%v The binary %v is not currently installed", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

// NoBinariesInstalledError occurs if you try to operate on a binary but no binaries are installed.
type NoBinariesInstalledError struct{}

func (e NoBinariesInstalledError) Error() string {
	return fmt.Sprintf("%v No binaries are currently installed", constants.RedColor("Error:"))
}

// UnrecognizedInputError occurs if the input is not recognized as a URL or GitHub repo.
type UnrecognizedInputError struct {
	Input string
}

func (e UnrecognizedInputError) Error() string {
	return fmt.Sprintf("%v Input was not recognized as a URL or GitHub repo: %v",
		constants.RedColor("Error:"), constants.RedColor(e.Input))
}

// InstalledFromURLError occurs if you try to perform GitHub specific actions on a binary installed directly from a URL.
type InstalledFromURLError struct {
	Binary string
}

func (e InstalledFromURLError) Error() string {
	return fmt.Sprintf("%v The %v binary was installed directly from a URL", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

// AlreadyInstalledLatestTagError occurs if you try to upgrade a binary but the latest version is already installed.
type AlreadyInstalledLatestTagError struct {
	Tag string
}

func (e AlreadyInstalledLatestTagError) Error() string {
	return fmt.Sprintf("%v The latest tag %v is already installed", constants.RedColor("Error:"), constants.RedColor(e.Tag))
}

// NoGithubSearchResultsError occurs if the GitHub search API returns no items.
type NoGithubSearchResultsError struct {
	SearchQuery string
}

func (e NoGithubSearchResultsError) Error() string {
	return fmt.Sprintf("%v No GitHub search results found for search query %v", constants.RedColor("Error:"), constants.RedColor(e.SearchQuery))
}

// InvalidGithubSearchQueryError occurs if the GitHub search query contains invalid characters.
type InvalidGithubSearchQueryError struct {
	Query string
}

func (e InvalidGithubSearchQueryError) Error() string {
	return fmt.Sprintf(`%v The search query "%v" contains invalid characters`,
		constants.RedColor("Error:"), constants.RedColor(e.Query))
}
