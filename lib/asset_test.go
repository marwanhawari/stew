package stew

import (
	"reflect"
	"testing"
)

// The asset lists in these tests are from real GitHub releases
func Test_rankAssets(t *testing.T) {
	tests := []struct {
		name     string
		userOS   string
		userArch string
		repo     string
		assets   []string
		want     []string
	}{
		{
			name:     "go names",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "age",
			assets:   []string{"age-v1.3.2-darwin-amd64.tar.gz", "age-v1.3.2-darwin-arm64.tar.gz", "age-v1.3.2-freebsd-amd64.tar.gz", "age-v1.3.2-linux-amd64.tar.gz", "age-v1.3.2-linux-arm.tar.gz", "age-v1.3.2-linux-arm64.tar.gz", "age-v1.3.2-windows-amd64.zip"},
			want:     []string{"age-v1.3.2-linux-amd64.tar.gz"},
		},
		{
			name:     "arm is not arm64",
			userOS:   "linux",
			userArch: "arm64",
			repo:     "age",
			assets:   []string{"age-v1.3.2-linux-amd64.tar.gz", "age-v1.3.2-linux-arm.tar.gz", "age-v1.3.2-linux-arm64.tar.gz"},
			want:     []string{"age-v1.3.2-linux-arm64.tar.gz"},
		},
		{
			name:     "darwin is not windows",
			userOS:   "windows",
			userArch: "amd64",
			repo:     "age",
			assets:   []string{"age-v1.3.2-darwin-amd64.tar.gz", "age-v1.3.2-linux-amd64.tar.gz", "age-v1.3.2-windows-amd64.zip"},
			want:     []string{"age-v1.3.2-windows-amd64.zip"},
		},
		{
			name:     "rust target triple prefers gnu on linux",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "starship",
			assets:   []string{"starship-aarch64-unknown-linux-musl.tar.gz", "starship-x86_64-apple-darwin.tar.gz", "starship-x86_64-pc-windows-msvc.zip", "starship-x86_64-unknown-linux-gnu.tar.gz", "starship-x86_64-unknown-linux-musl.tar.gz"},
			want:     []string{"starship-x86_64-unknown-linux-gnu.tar.gz"},
		},
		{
			name:     "musl when there is no other libc",
			userOS:   "linux",
			userArch: "arm64",
			repo:     "starship",
			assets:   []string{"starship-aarch64-unknown-linux-musl.tar.gz", "starship-x86_64-unknown-linux-gnu.tar.gz", "starship-x86_64-unknown-linux-musl.tar.gz"},
			want:     []string{"starship-aarch64-unknown-linux-musl.tar.gz"},
		},
		{
			name:     "msvc is preferred on windows",
			userOS:   "windows",
			userArch: "amd64",
			repo:     "bottom",
			assets:   []string{"bottom_x86_64-pc-windows-gnu.zip", "bottom_x86_64-pc-windows-msvc.zip", "bottom_i686-pc-windows-msvc.zip", "bottom_x86_64-unknown-linux-gnu.tar.gz"},
			want:     []string{"bottom_x86_64-pc-windows-msvc.zip"},
		},
		{
			name:     "linux64 and a build variant",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "micro",
			assets:   []string{"micro-2.0.15-freebsd64.tar.gz", "micro-2.0.15-linux-arm.tar.gz", "micro-2.0.15-linux-arm64.tar.gz", "micro-2.0.15-linux32.tar.gz", "micro-2.0.15-linux64-static.tar.gz", "micro-2.0.15-linux64.tar.gz", "micro-2.0.15-macos-arm64.tar.gz", "micro-2.0.15-osx.tar.gz", "micro-2.0.15-win64.zip"},
			want:     []string{"micro-2.0.15-linux64.tar.gz"},
		},
		{
			name:     "asset without an arch when no asset has the arch",
			userOS:   "darwin",
			userArch: "amd64",
			repo:     "micro",
			assets:   []string{"micro-2.0.15-linux64.tar.gz", "micro-2.0.15-macos-arm64.tar.gz", "micro-2.0.15-osx.tar.gz", "micro-2.0.15-win64.zip"},
			want:     []string{"micro-2.0.15-osx.tar.gz"},
		},
		{
			name:     "darwin arm64 uses amd64 when there is no arm64 asset",
			userOS:   "darwin",
			userArch: "arm64",
			repo:     "tool",
			assets:   []string{"tool-darwin-amd64.tar.gz", "tool-linux-arm64.tar.gz", "tool-linux-amd64.tar.gz"},
			want:     []string{"tool-darwin-amd64.tar.gz"},
		},
		{
			name:     "darwin universal",
			userOS:   "darwin",
			userArch: "arm64",
			repo:     "tool",
			assets:   []string{"tool_macos_universal.tar.gz", "tool_linux_arm64.tar.gz"},
			want:     []string{"tool_macos_universal.tar.gz"},
		},
		{
			name:     "binary is preferred to an archive",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "mise",
			assets:   []string{"mise-v2026.10.0-linux-arm64", "mise-v2026.10.0-linux-x64", "mise-v2026.10.0-linux-x64-musl", "mise-v2026.10.0-linux-x64-musl.tar.gz", "mise-v2026.10.0-linux-x64.tar.gz", "mise-v2026.10.0-linux-x64.tar.xz", "mise-v2026.10.0-macos-x64", "mise-v2026.10.0-windows-x64.exe"},
			want:     []string{"mise-v2026.10.0-linux-x64"},
		},
		{
			name:     "tar.xz is preferred to tar.gz",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "shellcheck",
			assets:   []string{"shellcheck-v0.11.0.darwin.x86_64.tar.gz", "shellcheck-v0.11.0.linux.aarch64.tar.gz", "shellcheck-v0.11.0.linux.x86_64.tar.gz", "shellcheck-v0.11.0.linux.x86_64.tar.xz"},
			want:     []string{"shellcheck-v0.11.0.linux.x86_64.tar.xz"},
		},
		{
			name:     "same archive in two formats",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "micro",
			assets:   []string{"micro-2.0.14-linux64.tgz", "micro-2.0.14-linux64.tar.gz", "micro-2.0.14-linux64-static.tgz", "micro-2.0.14-linux-arm64.tar.gz"},
			want:     []string{"micro-2.0.14-linux64.tar.gz"},
		},
		{
			name:     "formats in the order of the typical download size",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "tool",
			assets:   []string{"tool-linux-amd64.tar", "tool-linux-amd64.tar.lz4", "tool-linux-amd64.zip", "tool-linux-amd64.tar.gz", "tool-linux-amd64.tar.bz2", "tool-linux-amd64.tar.br", "tool-linux-amd64.tar.zst", "tool-linux-amd64.tar.xz"},
			want:     []string{"tool-linux-amd64.tar.xz"},
		},
		{
			name:     "zstd is preferred to gzip",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "tool",
			assets:   []string{"tool-linux-amd64.tar.gz", "tool-linux-amd64.tar.zst", "tool-linux-amd64.zip"},
			want:     []string{"tool-linux-amd64.tar.zst"},
		},
		{
			name:     "exe is preferred to zip on windows",
			userOS:   "windows",
			userArch: "amd64",
			repo:     "tool",
			assets:   []string{"tool_windows_amd64.exe", "tool_windows_amd64.zip", "tool_linux_amd64.tar.gz"},
			want:     []string{"tool_windows_amd64.exe"},
		},
		{
			name:     "program with the shortest name",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "zellij",
			assets:   []string{"zellij-aarch64-unknown-linux-musl.tar.gz", "zellij-no-web-x86_64-unknown-linux-musl.tar.gz", "zellij-x86_64-apple-darwin.tar.gz", "zellij-x86_64-unknown-linux-musl.tar.gz"},
			want:     []string{"zellij-x86_64-unknown-linux-musl.tar.gz"},
		},
		{
			name:     "program with the name of the repo",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "deno",
			assets:   []string{"denort-x86_64-unknown-linux-gnu.zip", "deno-aarch64-unknown-linux-gnu.zip", "deno-x86_64-unknown-linux-gnu.zip", "libdenort-x86_64-unknown-linux-gnu.zip"},
			want:     []string{"deno-x86_64-unknown-linux-gnu.zip"},
		},
		{
			name:     "linux android is not linux",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "bun",
			assets:   []string{"bun-linux-aarch64.zip", "bun-linux-x64-android.zip", "bun-linux-x64-baseline.zip", "bun-linux-x64-musl.zip", "bun-linux-x64-profile.zip", "bun-linux-x64.zip", "bun-windows-x64.zip"},
			want:     []string{"bun-linux-x64.zip"},
		},
		{
			name:     "exe without the windows name",
			userOS:   "windows",
			userArch: "amd64",
			repo:     "yt-dlp",
			assets:   []string{"yt-dlp", "yt-dlp.exe", "yt-dlp_arm64.exe", "yt-dlp_linux", "yt-dlp_macos", "yt-dlp_x86.exe"},
			want:     []string{"yt-dlp.exe"},
		},
		{
			name:     "one program for all operating systems",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "nextflow",
			assets:   []string{"nextflow", "nextflow-25.10.8-dist"},
			want:     []string{"nextflow"},
		},
		{
			name:     "archive without an OS is not selected",
			userOS:   "darwin",
			userArch: "arm64",
			repo:     "tool",
			assets:   []string{"tool-1.0.0.tar.gz", "tool-linux-amd64.tar.gz"},
			want:     nil,
		},
		{
			name:     "asset without an arch is not selected for linux arm64",
			userOS:   "linux",
			userArch: "arm64",
			repo:     "tool",
			assets:   []string{"tool-linux.tar.gz", "tool-darwin-arm64.tar.gz"},
			want:     nil,
		},
		{
			name:     "two programs in one release",
			userOS:   "linux",
			userArch: "amd64",
			repo:     "fuc",
			assets:   []string{"aarch64-unknown-linux-gnu-cpz", "x86_64-unknown-linux-gnu-cpz", "x86_64-unknown-linux-gnu-rmz", "x86_64-unknown-linux-musl-cpz", "x86_64-unknown-linux-musl-rmz"},
			want:     []string{"x86_64-unknown-linux-gnu-cpz", "x86_64-unknown-linux-gnu-rmz"},
		},
		{
			name:     "no asset for the arch",
			userOS:   "linux",
			userArch: "arm64",
			repo:     "cloudctl",
			assets:   []string{"cloudctl-darwin-arm64", "cloudctl-linux-amd64", "cloudctl-windows-amd64"},
			want:     nil,
		},
		{
			name:     "OS and arch that have no pattern",
			userOS:   "freebsd",
			userArch: "riscv64",
			repo:     "tool",
			assets:   []string{"tool-freebsd-amd64.tar.gz", "tool-freebsd-riscv64.tar.gz", "tool-linux-riscv64.tar.gz"},
			want:     []string{"tool-freebsd-riscv64.tar.gz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := rankAssets(tt.userOS, tt.userArch, tt.repo, tt.assets)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rankAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rankAssets_manualSelection(t *testing.T) {
	tests := []struct {
		name     string
		userOS   string
		userArch string
		assets   []string
		want     []string
	}{
		{
			name:     "assets for the OS and arch with the best first",
			userOS:   "linux",
			userArch: "amd64",
			assets:   []string{"tool-b-linux-amd64-musl.tar.gz", "tool-a-linux-amd64.tar.gz", "tool-b-linux-amd64.tar.gz", "tool-a-darwin-amd64.tar.gz", "tool-a-linux-arm64.tar.gz"},
			want:     []string{"tool-a-linux-amd64.tar.gz", "tool-b-linux-amd64.tar.gz", "tool-b-linux-amd64-musl.tar.gz"},
		},
		{
			name:     "assets for the OS when no asset is for the arch",
			userOS:   "linux",
			userArch: "arm64",
			assets:   []string{"tool-linux-amd64.tar.gz", "tool-linux-386.tar.gz", "tool-darwin-arm64.tar.gz"},
			want:     []string{"tool-linux-amd64.tar.gz", "tool-linux-386.tar.gz"},
		},
		{
			name:     "all assets when no asset is for the OS",
			userOS:   "darwin",
			userArch: "arm64",
			assets:   []string{"tool-linux-amd64.tar.gz", "tool-windows-amd64.zip"},
			want:     []string{"tool-linux-amd64.tar.gz", "tool-windows-amd64.zip"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := rankAssets(tt.userOS, tt.userArch, "tool", tt.assets)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rankAssets() manual selection = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterReleaseAssets_metadata(t *testing.T) {
	assets := []string{
		"tool-linux-amd64.tar.gz",
		"tool-linux-amd64.tar.gz.sha256",
		"tool-linux-amd64.tar.gz.sha",
		"tool-linux-amd64.tar.gz.md5",
		"tool-linux-amd64.tar.gz.sig",
		"tool-linux-amd64.tar.gz.asc",
		"tool-linux-amd64.tar.gz.minisig",
		"tool-linux-amd64.tar.gz.proof",
		"tool-linux-amd64.tar.gz.sbom.json",
		"tool-linux-amd64.tar.gz.zsync",
		"tool-linux-amd64.from-1.0.0.bsdiff",
		"tool-linux-amd64.tar.lz",
		"tool-windows-amd64.7z",
		"tool-1.0.0-py3-none-any.whl",
		"tool_1.0.0_checksums.txt",
		"checksums",
		"SHASUMS256.txt",
		"SHA2-256SUMS",
		"B3SUMS",
		"dist-manifest.json",
		"source.tar.gz",
		"tool_src.tar.gz",
		"tool.1",
		"tool",
	}
	want := []string{"tool-linux-amd64.tar.gz", "tool"}
	if got := filterReleaseAssets(assets); !reflect.DeepEqual(got, want) {
		t.Errorf("filterReleaseAssets() = %v, want %v", got, want)
	}
}

func TestMatchPreviousAsset(t *testing.T) {
	tests := []struct {
		name          string
		previousAsset string
		previousTag   string
		newTag        string
		releaseAssets []string
		want          string
		wantFound     bool
	}{
		{
			name:          "tag in the asset name",
			previousAsset: "mise-v2026.9.18-linux-x64",
			previousTag:   "v2026.9.18",
			newTag:        "v2026.10.0",
			releaseAssets: []string{"mise-v2026.10.0-linux-arm64", "mise-v2026.10.0-linux-x64", "mise-v2026.10.0-linux-x64-musl", "mise-v2026.10.0-linux-x64.tar.gz"},
			want:          "mise-v2026.10.0-linux-x64",
			wantFound:     true,
		},
		{
			name:          "version without the v of the tag",
			previousAsset: "gh_2.101.0_linux_amd64.tar.gz",
			previousTag:   "v2.101.0",
			newTag:        "v2.102.0",
			releaseAssets: []string{"gh_2.102.0_linux_amd64.tar.gz", "gh_2.102.0_linux_arm64.tar.gz", "gh_2.102.0_linux_amd64.deb"},
			want:          "gh_2.102.0_linux_amd64.tar.gz",
			wantFound:     true,
		},
		{
			name:          "tag with a prefix",
			previousAsset: "bun-1.2.0-linux-x64-baseline.zip",
			previousTag:   "bun-v1.2.0",
			newTag:        "bun-v1.3.0",
			releaseAssets: []string{"bun-1.3.0-linux-x64.zip", "bun-1.3.0-linux-x64-baseline.zip"},
			want:          "bun-1.3.0-linux-x64-baseline.zip",
			wantFound:     true,
		},
		{
			name:          "no version in the asset name",
			previousAsset: "ffmpeg-master-latest-linux64-gpl.tar.xz",
			previousTag:   "autobuild-2026-09-01-12-00",
			newTag:        "autobuild-2026-10-01-12-00",
			releaseAssets: []string{"ffmpeg-master-latest-linux64-gpl-shared.tar.xz", "ffmpeg-master-latest-linux64-gpl.tar.xz", "ffmpeg-master-latest-linux64-lgpl.tar.xz"},
			want:          "ffmpeg-master-latest-linux64-gpl.tar.xz",
			wantFound:     true,
		},
		{
			name:          "the asset name changed",
			previousAsset: "rig-linux-0.8.0.tar.gz",
			previousTag:   "v0.8.0",
			newTag:        "v0.10.0",
			releaseAssets: []string{"rig-linux-aarch64-0.10.0.tar.gz", "rig-linux-x86_64-0.10.0.tar.gz"},
			want:          "",
			wantFound:     false,
		},
		{
			name:          "no previous asset",
			previousAsset: "",
			previousTag:   "v1.0.0",
			newTag:        "v1.1.0",
			releaseAssets: []string{"tool-linux-amd64.tar.gz"},
			want:          "",
			wantFound:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotFound := MatchPreviousAsset(tt.previousAsset, tt.previousTag, tt.newTag, tt.releaseAssets)
			if got != tt.want || gotFound != tt.wantFound {
				t.Errorf("MatchPreviousAsset() = %v, %v, want %v, %v", got, gotFound, tt.want, tt.wantFound)
			}
		})
	}
}
