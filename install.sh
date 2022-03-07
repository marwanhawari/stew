#!/usr/bin/env bash

# This install script does 3 things:
# 1. Create the stew directory structure
# 2. Download the stew binary
# 3. Add the stew binary path to PATH in ~/.zshrc or ~/.bashrc

os=""
arch=""
exe=""
defaultStewPath=""
configPath=""

# Detect os
case "$(uname -s)" in
    Darwin)
        os="darwin"

        if [ -z "$XDG_DATA_HOME" ]
        then
            defaultStewPath="$HOME/.local/share/stew"
        else
            defaultStewPath="$XDG_DATA_HOME/stew"
        fi

        if [ -z "$XDG_CONFIG_HOME" ]
        then
            configPath="$HOME/.config/stew"
        else
            configPath="$XDG_CONFIG_HOME/stew"
        fi
        ;;
    Linux)
        os="linux"
        if [ -z "$XDG_DATA_HOME" ]
        then
            defaultStewPath="$HOME/.local/share/stew"
        else
            defaultStewPath="$XDG_DATA_HOME/stew"
        fi

        if [ -z "$XDG_CONFIG_HOME" ]
        then
            configPath="$HOME/.config/stew"
        else
            configPath="$XDG_CONFIG_HOME/stew"
        fi
        ;;
    CYGWIN*|MSYS*|MINGW*)
        os="windows"
        exe=".exe"
        defaultStewPath="$HOME/AppData/Local/stew"
        configPath="$HOME/AppData/Local/stew/Config"
        ;;
esac

# Detect arch
case "$(uname -m)" in
    amd64|x86_64)
        arch="amd64"
        ;;
    arm64|aarch64)
        arch="arm64"
        ;;
    *386)
        arch="386"
        ;;
esac

if [ "$os" = "" ] || [ "$arch" = "" ]
then
    echo ""
    echo "|||||||||||||||||||||"
    echo "||      Error      ||"
    echo "|||||||||||||||||||||"
    echo ""
    echo "Your current OS/arch is not supported by stew"
    echo ""
    exit 1
fi

# 1. Create the stew directory structure
stewPath=""
stewBinPath=""

read -r -t 60 -p "Set the stewPath. This will contain all stew data other than the binaries. (${defaultStewPath}): " stewPathInput
if [ -z "$stewPathInput" ]
then
    stewPath="${defaultStewPath}"
else
    stewPath="${stewPathInput/#~/$HOME}"
    stewPath="${stewPath/#\$HOME/$HOME}"
    stewPath="${stewPath/#\$PWD/$PWD}"
    if [ -x "$(command -v dirname)" ] && [ -x "$(command -v basename)" ]
    then
        stewPath="$(cd "$(dirname "$stewPath")" || exit; pwd)/$(basename "$stewPath")"
    fi
fi

read -r -t 60 -p "Set the stewBinPath. This is where the binaries will be installed by stew. (${defaultStewPath}/bin): " stewBinPathInput
if [ -z "$stewBinPathInput" ]
then
    stewBinPath="${defaultStewPath}/bin"
else
    stewBinPath="${stewBinPathInput/#~/$HOME}"
    stewBinPath="${stewBinPath/#\$HOME/$HOME}"
    stewBinPath="${stewBinPath/#\$PWD/$PWD}"
    if [ -x "$(command -v dirname)" ] && [ -x "$(command -v basename)" ]
    then
        stewBinPath="$(cd "$(dirname "$stewBinPath")" || exit; pwd)/$(basename "$stewBinPath")"
    fi
fi

mkdir -p "${stewPath}/bin"
mkdir -p "${stewPath}/pkg"
mkdir -p "${stewBinPath}"
mkdir -p "${configPath}"

echo "{
	\"stewPath\": \"${stewPath}\",
	\"stewBinPath\": \"${stewBinPath}\"
}" > "${configPath}/config.json"

# 2. Download the stew binary
curl -o "${stewBinPath}/stew${exe}" -fsSL https://github.com/marwanhawari/stew/releases/latest/download/stew-${os}-${arch}${exe}
chmod +x "${stewBinPath}/stew${exe}"

# 3. Add the stew binary path to PATH in ~/.zshrc or ~/.bashrc
if [ -f "$HOME"/.zshrc ]
then
    echo "export PATH=\"${stewBinPath}:\$PATH\"" >> "$HOME"/.zshrc
elif [ -f "$HOME"/.bashrc ]
then
    echo "export PATH=\"${stewBinPath}:\$PATH\"" >> "$HOME"/.bashrc
else
    echo "Make sure to add ${stewBinPath} to PATH"
fi

echo ""
echo "|||||||||||||||||||||"
echo "||     Success     ||"
echo "|||||||||||||||||||||"
echo ""
echo "Start a new terminal session to start using stew"
echo ""
