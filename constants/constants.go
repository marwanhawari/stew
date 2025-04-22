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

// CyanColor makes text cyan
var CyanColor = color.New(color.FgCyan, color.OpBold).Render

// BoldColor makes text bold
var BoldColor = color.New(color.OpBold).Render

// LoadingSpinner is a reusable loading spinner
var LoadingSpinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithColor("cyan"), spinner.WithHiddenCursor(true))

// RegexDarwin is a regular express for darwin systems
var RegexDarwin = `(?i)(darwin|mac(os)?|apple|osx)`

// RegexWindows is a regular express for windows systems
var RegexWindows = `(?i)(windows|win|.msi|.exe)`

// RegexArm64 is a regular express for arm64 architectures
var RegexArm64 = `(?i)(arm64|aarch64|arm64e)`

// RegexAmd64 is a regular express for amd64 architectures
var RegexAmd64 = `(?i)(x86_64|amd64|x64|amd64e)`

// Regex386 is a regular express for 386 architectures
var Regex386 = `(?i)(i?386|x86_32|amd32|x32)`

// RegexGithub is a regular express for valid GitHub repos
var RegexGithub = `(?i)^[A-Za-z0-9\-]+\/[A-Za-z0-9\_\.\-]+(@.+)?$`

// RegexGithubSearch is a regular express for valid GitHub search queries
var RegexGithubSearch = `(?i)^[A-Za-z0-9\_\.\-\/\:]+$`

// RegexURL is a regular express for valid URLs
var RegexURL = `(http|ftp|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`

// RegexChecksum is a regular expression for matching checksum files
var RegexChecksum = `\.(sha(256|512)(sum)?)$`

// StewOwner is the username of the stew github repo owner
var StewOwner = `marwanhawari`

// StewRepo is the name of the stew github repo
var StewRepo = `stew`
