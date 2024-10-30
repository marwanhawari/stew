#!/bin/sh

# This script installs Stew.
#
# Quick install: `curl https://github.com/marwanhawari/stew/install.sh | bash`
#
# Acknowledgments:
#   - getmic.ro: https://github.com/benweissmann/getmic.ro

set -e -u

githubLatestTag() {
  finalUrl=$(curl "https://github.com/$1/releases/latest" -s -L -I -o /dev/null -w '%{url_effective}')
  printf "%s\n" "${finalUrl##*v}"
}

platform=''
machine=$(uname -m)

if [ "${GETSTEW_PLATFORM:-x}" != "x" ]; then
  platform="$GETSTEW_PLATFORM"
else
  case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
    "linux")
      case "$machine" in
        "arm64"* | "aarch64"* ) platform='linux-arm64' ;;
        "arm"* | "aarch"*) platform='linux-arm' ;;
        *"86") platform='linux-386' ;;
        *"64") platform='linux-amd64' ;;
      esac
      ;;
    "darwin")
      case "$machine" in
        "arm64"* | "aarch64"* ) platform='darwin-arm64' ;;
        *"64") platform='darwin-amd64' ;;
      esac
      ;;
    "msys"*|"cygwin"*|"mingw"*|*"_nt"*|"win"*)
      case "$machine" in
        *"86") platform='windows-386' ;;
        *"64") platform='windows-amd64' ;;
      esac
      ;;
  esac
fi

if [ "x$platform" = "x" ]; then
  cat << 'EOM'
/=====================================\\
|      COULD NOT DETECT PLATFORM      |
\\=====================================/
Uh oh! We couldn't automatically detect your operating system.
To continue with installation, please choose from one of the following values:
- linux-arm
- linux-arm64
- linux-386
- linux-amd64
- darwin-amd64
- darwin-arm64
- windows-386
- windows-amd64
Export your selection as the GETSTEW_PLATFORM environment variable, and then
re-run this script.
For example:
  $ export GETSTEW_PLATFORM=linux-amd64
  $ curl https://github.com/marwanhawari/stew/install.sh | bash
EOM
  exit 1
else
  printf "Detected platform: %s\n" "$platform"
fi

TAG=v$(githubLatestTag marwanhawari/stew)

if [ "x$platform" = "xwindows_amd64" ] || [ "x$platform" = "xwindows_386" ]; then
  extension='zip'
else
  extension='tar.gz'
fi

printf "Latest Version: %s\n" "$TAG"
printf "Downloading https://github.com/marwanhawari/stew/releases/download/%s/stew-%s-%s.%s\n" "$TAG" "$TAG" "$platform" "$extension"

curl -L "https://github.com/marwanhawari/stew/releases/download/$TAG/stew-$TAG-$platform.$extension" > "stew.$extension"

case "$extension" in
  "zip") unzip -j "stew.$extension" -d "stew" ;;
  "tar.gz") tar -xvzf "stew.$extension" "stew" ;;
esac

chmod 777 stew

cat <<-'EOM'
stew has been downloaded to the current directory.
You can run it with:
./stew
EOM
