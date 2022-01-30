package stew

import (
	"fmt"

	"github.com/marwanhawari/stew/constants"
)

type NonZeroStatusCodeError struct {
	StatusCode int
}

func (e NonZeroStatusCodeError) Error() string {
	return fmt.Sprintf("%v Received non-zero status code from HTTP request: %v", constants.RedColor("Error:"), constants.YellowColor(e.StatusCode))
}

type ReleasesNotFoundError struct {
	Owner string
	Repo  string
}

func (e ReleasesNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find any releases for %v", constants.RedColor("Error:"), constants.YellowColor("https://github.com/"+e.Owner+"/"+e.Repo))
}

type AssetsNotFoundError struct {
	Tag string
}

func (e AssetsNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find any assets for release %v", constants.RedColor("Error:"), constants.YellowColor(e.Tag))
}

type NoPackagesInLockfileError struct {
}

func (e NoPackagesInLockfileError) Error() string {
	return fmt.Sprintf("%v Cannot remove from an empty packages slice in the lockfile", constants.RedColor("Error:"))
}

type IndexOutOfBoundsInLockfileError struct {
}

func (e IndexOutOfBoundsInLockfileError) Error() string {
	return fmt.Sprintf("%v Index out of bounds in lockfile packages", constants.RedColor("Error:"))
}

type ExitUserSelectionError struct {
	Err error
}

func (e ExitUserSelectionError) Error() string {
	return fmt.Sprintf("%v Exited from user selection: %v", constants.RedColor("Error:"), constants.YellowColor(e.Err))
}

type StewPathNotFoundError struct {
	StewPath string
}

func (e StewPathNotFoundError) Error() string {
	return fmt.Sprintf("%v Could not find the stew path at %v", constants.RedColor("Error:"), e.StewPath)
}

type NonZeroStatusCodeDownloadError struct {
	StatusCode int
}

func (e NonZeroStatusCodeDownloadError) Error() string {
	return fmt.Sprintf("%v Received non-zero status code from HTTP request when attempting to download a file: %v", constants.RedColor("Error:"), constants.YellowColor(e.StatusCode))
}

type EmptyCLIInputError struct {
}

func (e EmptyCLIInputError) Error() string {
	return fmt.Sprintf("%v Input cannot be empty. Use the --help flag for more info", constants.RedColor("Error:"))
}

type UninstallCLIInputError struct {
}

func (e UninstallCLIInputError) Error() string {
	return fmt.Sprintf("%v Cannot use the --all flag with a positional argument", constants.RedColor("Error:"))
}

type AssetAlreadyDownloadedError struct {
	Asset string
}

func (e AssetAlreadyDownloadedError) Error() string {
	return fmt.Sprintf("%v The %v asset has already been downloaded and installed", constants.RedColor("Error:"), constants.RedColor(e.Asset))
}

type AbortBinaryOverwriteError struct {
	Binary string
}

func (e AbortBinaryOverwriteError) Error() string {
	return fmt.Sprintf("%v Overwrite of %v aborted", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

type BinaryNotInstalledError struct {
	Binary string
}

func (e BinaryNotInstalledError) Error() string {
	return fmt.Sprintf("%v The binary %v is not currently installed", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

type NoBinariesInstalledError struct {
}

func (e NoBinariesInstalledError) Error() string {
	return fmt.Sprintf("%v No binaries are currently installed", constants.RedColor("Error:"))
}

type UnrecognizedInputError struct {
}

func (e UnrecognizedInputError) Error() string {
	return fmt.Sprintf("%v Input was not recognized as a URL or GitHub repo", constants.RedColor("Error:"))
}

type InstalledFromURLError struct {
	Binary string
}

func (e InstalledFromURLError) Error() string {
	return fmt.Sprintf("%v The %v binary was installed directly from a URL", constants.RedColor("Error:"), constants.RedColor(e.Binary))
}

type AlreadyInstalledLatestTagError struct {
	Tag string
}

func (e AlreadyInstalledLatestTagError) Error() string {
	return fmt.Sprintf("%v The latest tag %v is already installed", constants.RedColor("Error:"), constants.RedColor(e.Tag))
}
