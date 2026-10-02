package stew

import (
	"regexp"
	"sort"
	"strings"

	"github.com/marwanhawari/stew/constants"
)

// assetPattern finds one class of OS, arch, or libc in an asset name.
// The matched text is replaced so that a later pattern cannot match it again. For example, arm must not match in arm64.
type assetPattern struct {
	class       string
	re          *regexp.Regexp
	replacement string
}

const (
	archUniversal = "universal"
	classOther    = "other"

	libcGnu  = "gnu"
	libcMusl = "musl"
	libcMsvc = "msvc"
)

// The order is important. The first patterns remove text that the subsequent patterns must not see.
var archPatterns = []assetPattern{
	{"amd64", regexp.MustCompile(`(linux|win|windows|macos|osx|darwin|bsd|illumos|solaris)64`), "${1}-"},
	{"386", regexp.MustCompile(`(linux|win|windows|macos|osx|darwin|bsd|illumos|solaris)32`), "${1}-"},
	{"amd64", regexp.MustCompile(`x86[_-]64|amd64|(^|[^a-z0-9])x64|(^|[^a-z0-9])64[_-]?bit`), "-"},
	{"arm64", regexp.MustCompile(`arm64e?|aarch64`), "-"},
	{classOther, regexp.MustCompile(`riscv64(gc)?|ppc64(le|el)?|powerpc(64)?(le)?|s390x|mips(64)?(r6)?(el|le)?|loong(arch)?64|m68k|sparc(64)?|wasm(32|64)?`), "-"},
	{"386", regexp.MustCompile(`i[3-6]86|x86[_-]32|x86|(^|[^a-z0-9])386|(^|[^a-z0-9])32[_-]?bit`), "-"},
	{"arm", regexp.MustCompile(`(^|[^a-z])arm(v[5-8][a-z]*|hf|el|32)?`), "-"},
	{archUniversal, regexp.MustCompile(`universal2?`), "-"},
}

var osPatterns = []assetPattern{
	{"darwin", regexp.MustCompile(`darwin|macos|apple|(^|[^a-z])(mac|osx)([^a-z]|$)`), "-"},
	{"windows", regexp.MustCompile(`windows|(^|[^a-z])win([^a-z]|$)`), "-"},
	{"linux", regexp.MustCompile(`linux`), "-"},
	{classOther, regexp.MustCompile(`freebsd|netbsd|openbsd|dragonfly|solaris|illumos|sunos|android|plan9|haiku|(^|[^a-z])(bsd|aix|wasi)([^a-z]|$)`), "-"},
}

var libcPatterns = []assetPattern{
	{libcGnu, regexp.MustCompile(`(^|[^a-z])(gnu|glibc)(eabi(hf)?)?`), "-"},
	{libcMusl, regexp.MustCompile(`(^|[^a-z])musl(eabi(hf)?)?`), "-"},
	{libcMsvc, regexp.MustCompile(`(^|[^a-z])msvc`), "-"},
}

// formatRankBinary is the rank of a binary that is not in an archive and is not compressed
const formatRankBinary = 0

// assetFormats gives the extensions that stew can install in the order of preference.
// A binary that is not in an archive is first. It is the largest download, but it is always one complete program.
// An archive can contain a program that needs the other files in the archive.
// The compressed formats are in the order from the smallest typical download size to the largest: xz, zstd, brotli, bzip2, gzip, zip, lz4, snappy, no compression.
// Each extension has a different rank, thus a release with the same build in two formats has one best asset.
// A longer extension must be before the shorter extension that it ends with.
var assetFormats = []struct {
	extension string
	rank      int
}{
	{".tar.xz", 1}, {".txz", 2}, {".xz", 3},
	{".tar.zst", 4}, {".tzst", 5}, {".zst", 6},
	{".tar.br", 7}, {".tbr", 8}, {".br", 9},
	{".tar.bz2", 10}, {".tbz2", 11}, {".tbz", 12}, {".bz2", 13},
	{".tar.gz", 14}, {".tgz", 15}, {".gz", 16},
	{".zip", 17}, {".rar", 18},
	{".tar.lz4", 19}, {".tlz4", 20}, {".lz4", 21},
	{".tar.sz", 22}, {".tsz", 23}, {".sz", 24},
	{".tar", 25},
	{".exe", formatRankBinary}, {".appimage", formatRankBinary}, {".bin", formatRankBinary},
}

var (
	reMetadata          = regexp.MustCompile(constants.RegexMetadata)
	reUnsupportedFormat = regexp.MustCompile(constants.RegexUnsupportedFormat)
	reNotAlphanumeric   = regexp.MustCompile(`[^a-z0-9]+`)
	reVersionToken      = regexp.MustCompile(`^v?[0-9]`)
	reUnknownExtension  = regexp.MustCompile(`\.[a-z][a-z0-9]*$`)
	// reVendorWord matches the vendor part of a target triple, for example x86_64-unknown-linux-gnu
	reVendorWord = regexp.MustCompile(`^(unknown|pc)$`)
)

// assetInfo is the data that stew reads from the name of a release asset
type assetInfo struct {
	name       string
	oses       map[string]bool
	archs      map[string]bool
	libc       string
	formatRank int
	// unknownExtension is true if the name has an extension that is not in assetFormats
	unknownExtension bool
	// words are the parts of the name that are not an OS, an arch, a libc, a version, or the extension.
	// These usually are the program name and the names of build variants.
	words []string
}

func applyPatterns(name string, patterns []assetPattern) (string, map[string]bool) {
	classes := map[string]bool{}
	for _, pattern := range patterns {
		if pattern.re.MatchString(name) {
			classes[pattern.class] = true
			name = pattern.re.ReplaceAllString(name, pattern.replacement)
		}
	}
	return name, classes
}

func isKnownClass(class string, patterns []assetPattern) bool {
	for _, pattern := range patterns {
		if pattern.class == class {
			return true
		}
	}
	return false
}

// parseAsset reads the OS, the arch, the libc, and the format from the name of a release asset
func parseAsset(asset, userOS, userArch string) assetInfo {
	info := assetInfo{name: asset, formatRank: formatRankBinary}
	name := strings.ToLower(asset)

	for _, format := range assetFormats {
		if strings.HasSuffix(name, format.extension) {
			info.formatRank = format.rank
			name = strings.TrimSuffix(name, format.extension)
			break
		}
	}
	info.unknownExtension = name == strings.ToLower(asset) && reUnknownExtension.MatchString(name)

	// An OS or an arch that stew has no pattern for is found by its Go name
	unknownUserOS := !isKnownClass(userOS, osPatterns) && strings.Contains(name, strings.ToLower(userOS))
	if unknownUserOS {
		name = strings.ReplaceAll(name, strings.ToLower(userOS), "-")
	}
	unknownUserArch := !isKnownClass(userArch, archPatterns) && strings.Contains(name, strings.ToLower(userArch))
	if unknownUserArch {
		name = strings.ReplaceAll(name, strings.ToLower(userArch), "-")
	}

	name, info.archs = applyPatterns(name, archPatterns)
	name, info.oses = applyPatterns(name, osPatterns)
	if strings.HasSuffix(strings.ToLower(asset), ".exe") {
		info.oses["windows"] = true
	}
	if unknownUserOS {
		info.oses[userOS] = true
	}
	if unknownUserArch {
		info.archs[userArch] = true
	}

	name, libcs := applyPatterns(name, libcPatterns)
	for _, libc := range []string{libcMusl, libcGnu, libcMsvc} {
		if libcs[libc] {
			info.libc = libc
			break
		}
	}

	for _, word := range reNotAlphanumeric.Split(name, -1) {
		if word == "" || reVersionToken.MatchString(word) || reVendorWord.MatchString(word) {
			continue
		}
		info.words = append(info.words, word)
	}

	return info
}

// assetRank is used to compare assets that match the OS and the arch. A smaller value is better.
type assetRank struct {
	libc     int
	words    int
	repoName int
	format   int
}

func (info assetInfo) rank(userOS, repo string) assetRank {
	rank := assetRank{words: len(info.words), format: info.formatRank, repoName: 1}

	// The gnu libc is on almost all Linux systems. The msvc build is the native build on Windows.
	if userOS == "linux" && info.libc == libcMusl {
		rank.libc = 1
	}
	if userOS == "windows" && info.libc == libcGnu {
		rank.libc = 1
	}

	// A release can contain more than one program. Prefer the program that has the name of the repo.
	repoName := reNotAlphanumeric.ReplaceAllString(strings.ToLower(repo), "")
	if repoName != "" && strings.Join(info.words, "") == repoName {
		rank.repoName = 0
	}

	return rank
}

func (rank assetRank) less(other assetRank) bool {
	if rank.libc != other.libc {
		return rank.libc < other.libc
	}
	if rank.words != other.words {
		return rank.words < other.words
	}
	if rank.repoName != other.repoName {
		return rank.repoName < other.repoName
	}
	return rank.format < other.format
}

func filterAssetInfos(infos []assetInfo, keep func(assetInfo) bool) []assetInfo {
	var kept []assetInfo
	for _, info := range infos {
		if keep(info) {
			kept = append(kept, info)
		}
	}
	return kept
}

func assetNames(infos []assetInfo) []string {
	var names []string
	for _, info := range infos {
		names = append(names, info.name)
	}
	return names
}

// rankAssets finds the release assets that match the OS and the arch and puts them in the order of preference.
// The first return value contains the best assets. If it contains exactly one asset, the detection was successful.
// The second return value contains all the assets that are applicable for a manual selection, with the best assets first.
func rankAssets(userOS, userArch, repo string, assets []string) ([]string, []string) {
	var infos []assetInfo
	for _, asset := range assets {
		infos = append(infos, parseAsset(asset, userOS, userArch))
	}

	// 1. The OS must match. An asset for two operating systems, for example linux-android, is not accepted.
	candidates := filterAssetInfos(infos, func(info assetInfo) bool {
		return info.oses[userOS] && len(info.oses) == 1
	})
	if len(candidates) == 0 {
		// Some releases contain one program for all operating systems, for example a script.
		// Archives without an OS in the name are not accepted, because they usually contain source code or documentation.
		candidates = filterAssetInfos(infos, func(info assetInfo) bool {
			return len(info.oses) == 0 && len(info.archs) == 0 && info.formatRank == formatRankBinary && !info.unknownExtension
		})
		if len(candidates) == 0 {
			return nil, assets
		}
		best := bestAssetInfos(candidates, userOS, repo)
		if len(best) == 1 {
			return assetNames(best), assetNames(candidates)
		}
		return nil, assets
	}
	osCandidates := candidates

	// 2. The arch must match
	candidates = filterAssetInfos(osCandidates, func(info assetInfo) bool {
		return info.archs[userArch] || (userOS == "darwin" && info.archs[archUniversal])
	})
	if len(candidates) == 0 && userOS == "darwin" && userArch == "arm64" {
		// Rosetta can run amd64 binaries
		candidates = filterAssetInfos(osCandidates, func(info assetInfo) bool { return info.archs["amd64"] })
	}
	if len(candidates) == 0 && (userArch == "amd64" || userOS == "darwin") {
		// An asset without an arch in its name is usually for amd64. On darwin it is usually a universal binary.
		candidates = filterAssetInfos(osCandidates, func(info assetInfo) bool { return len(info.archs) == 0 })
	}
	if len(candidates) == 0 {
		return nil, assetNames(osCandidates)
	}

	// 3. Put the assets in the order of preference
	best := bestAssetInfos(candidates, userOS, repo)

	return assetNames(best), assetNames(candidates)
}

// bestAssetInfos sorts the assets in the order of preference and returns the assets that have the best rank
func bestAssetInfos(candidates []assetInfo, userOS, repo string) []assetInfo {
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].rank(userOS, repo).less(candidates[j].rank(userOS, repo))
	})
	bestRank := candidates[0].rank(userOS, repo)
	return filterAssetInfos(candidates, func(info assetInfo) bool {
		return info.rank(userOS, repo) == bestRank
	})
}

// versionVariants gives the forms of a release tag that can be in the name of an asset. For example, the tag
// tool-v1.2.3 gives tool-v1.2.3, v1.2.3, and 1.2.3.
func versionVariants(tag string) []string {
	variants := []string{}
	if tag == "" {
		return variants
	}
	variants = append(variants, tag)
	if index := strings.IndexAny(tag, "0123456789"); index > 0 {
		if tag[index-1] == 'v' && index > 1 {
			variants = append(variants, tag[index-1:])
		}
		variants = append(variants, tag[index:])
	}
	return variants
}

func replaceVersion(asset, tag string) string {
	for _, variant := range versionVariants(tag) {
		if strings.Contains(asset, variant) {
			return strings.ReplaceAll(asset, variant, "\x00")
		}
	}
	return asset
}

// MatchPreviousAsset finds the asset of a new release that has the same name as the asset of the installed release,
// apart from the version. This keeps the selection that the user made at the installation.
func MatchPreviousAsset(previousAsset, previousTag, newTag string, releaseAssets []string) (string, bool) {
	if previousAsset == "" {
		return "", false
	}
	previousPattern := replaceVersion(previousAsset, previousTag)

	var matches []string
	for _, asset := range releaseAssets {
		if replaceVersion(asset, newTag) == previousPattern {
			matches = append(matches, asset)
		}
	}
	if len(matches) != 1 {
		return "", false
	}
	return matches[0], true
}
