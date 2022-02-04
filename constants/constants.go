package constants

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
)

var RedColor = color.New(color.FgRed, color.OpBold).Render
var GreenColor = color.New(color.FgGreen, color.OpBold).Render
var YellowColor = color.New(color.FgYellow, color.OpBold).Render
var BoldColor = color.New(color.OpBold).Render
var LoadingSpinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithColor("cyan"), spinner.WithHiddenCursor(true))

var RegexDarwin = `(?i)(darwin|mac(os)?|apple|osx)`
var RegexWindows = `(?i)(windows|win)`
var RegexArm64 = `(?i)(arm64|aarch64)`
var RegexAmd64 = `(?i)(x86_64|amd64|x64)`
var Regex386 = `(?i)(i?386|x86_32|amd32|x32)`

var RegexGithub = `(?i)^[A-Za-z0-9-]+\/[A-Za-z0-9_.-]+(@.+)?$`
var RegexGithubSearch = `(?i)^[A-Za-z0-9_.-]+$`
var RegexURL = `(http|ftp|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`
