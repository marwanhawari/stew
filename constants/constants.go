package constants

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
)

// RedColor makes text red
var RedColor = color.New(color.FgRed, color.OpBold).Render

// GreenColor makes text green
var GreenColor = color.New(color.FgGreen, color.OpBold).Render

// YellowColor makes text yellow
var YellowColor = color.New(color.FgYellow, color.OpBold).Render

// BoldColor makes text bold
var BoldColor = color.New(color.OpBold).Render

// LoadingSpinner is a reusable loading spinner
var LoadingSpinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithColor("cyan"), spinner.WithHiddenCursor(true))

// RegexGithub is a regular express for valid GitHub repos
var RegexGithub = `(?i)^[A-Za-z0-9\-]+\/[A-Za-z0-9\_\.\-]+(@.+)?$`

// RegexGithubSearch is a regular express for valid GitHub search queries
var RegexGithubSearch = `(?i)^[A-Za-z0-9\_\.\-\/\:]+$`

// RegexURL is a regular express for valid URLs
var RegexURL = `^(http|ftp|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`

// RegexChecksum is a regular expression for matching checksum files
var RegexChecksum = `\.(sha(256|512)(sum)?)$`

// RegexMetadata is a regular expression for matching assets that are not programs: checksums, signatures, attestations, update patches, source archives, man pages, and data files
var RegexMetadata = `(?i)(\.(sha[0-9]*(sum)?|md5|sig|asc|minisig|pem|crt|pub|proof|sbom|json|jsonl|attestation|bundle|zsync|bsdiff|txt|md|whl|ts|[1-9])$|checksum|(sha|md5|b3)[0-9-]*sums|(^|[-_.])(source|src)\.(tar|tgz|zip))`

// RegexUnsupportedFormat is a regular expression for matching archive formats that stew cannot extract
var RegexUnsupportedFormat = `(?i)\.(7z|lz|lzma|cab)$`

// RegexPackageInstaller is a regular expression for matching package-manager and installer assets
var RegexPackageInstaller = `(?i)(\.deb|\.rpm|\.apk|\.msi|\.pkg|\.dmg|\.pkg\.tar\.(gz|xz|zst))$`

// StewOwner is the username of the stew github repo owner
var StewOwner = `marwanhawari`

// StewRepo is the name of the stew github repo
var StewRepo = `stew`
