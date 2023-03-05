package constants

import (
	"regexp"

	"github.com/gookit/color"
)

// RedColor makes text red
var RedColor = color.New(color.FgRed, color.OpBold).Render

// GreenColor makes text green
var GreenColor = color.New(color.FgGreen, color.OpBold).Render

// YellowColor makes text yellow
var YellowColor = color.New(color.FgYellow, color.OpBold).Render

// RegexDarwin is a regular express for darwin systems
var RegexDarwin = regexp.MustCompile(`(?i)(darwin|mac(os)?|apple|osx)`)

// RegexWindows is a regular express for windows systems
var RegexWindows = regexp.MustCompile(`(?i)(windows|win)`)

// RegexArm64 is a regular express for arm64 architectures
var RegexArm64 = regexp.MustCompile(`(?i)(arm64|aarch64)`)

// RegexAmd64 is a regular express for amd64 architectures
var RegexAmd64 = regexp.MustCompile(`(?i)(x86_64|amd64|x64)`)

// Regex386 is a regular express for 386 architectures
var Regex386 = regexp.MustCompile(`(?i)(i?386|x86_32|amd32|x32)`)

// RegexGithub is a regular express for valid GitHub repos
var RegexGithub = regexp.MustCompile(`(?i)^([A-Za-z0-9\-]+)\/([A-Za-z0-9_.\-]+)(?:@([A-Za-z0-9_.\-]+))?(?:::([A-Za-z0-9_.\-]+))?(?:!!([A-Za-z0-9_.\-]+))?$`)

// RegexGithubSearch is a regular express for valid GitHub search queries
var RegexGithubSearch = regexp.MustCompile(`(?i)^[A-Za-z0-9\_\.\-\/]+$`)

// RegexURL is a regular express for valid URLs
var RegexURL = regexp.MustCompile(`(http|ftp|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`)
