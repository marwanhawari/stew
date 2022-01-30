#!/usr/bin/env sh

# This install script does 3 things:
# 1. Create the ~/.stew directory structure
# 2. Download the stew binary
# 3. Add ~/.stew/bin to ~/.zshrc and ~/.bashrc

os=""
arch=""
exe=""

# Detect os
case "$(uname -s)" in
    Darwin)
        os="darwin"
        ;;
    Linux)
        os="linux"
        ;;
    CYGWIN*|MSYS*|MINGW*)
        os="windows"
        exe=".exe"
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
    echo "\033[31m\033[1mError:\033[0m Your current OS/arch is not supported by stew"
    exit 1
fi

# 1. Create the ~/.stew directory structure
mkdir -p "$HOME"/.stew/bin
mkdir -p "$HOME"/.stew/pkg

# 2. Download the stew binary
curl -o "$HOME"/.stew/bin/stew${exe} -fsSL https://github.com/marwanhawari/stew/releases/latest/download/stew-${os}-${arch}${exe}
chmod +x "$HOME"/.stew/bin/stew${exe}

# 3. Add ~/.stew/bin to $PATH in ~/.zshrc or ~/.bashrc
if [ -f "$HOME"/.zshrc ]
then
    echo 'export PATH="$HOME/.stew/bin:$PATH"' >> "$HOME"/.zshrc
elif [ -f "$HOME"/.bashrc ]
then
    echo 'export PATH="$HOME/.stew/bin:$PATH"' >> "$HOME"/.bashrc
else
    echo "Make sure to add $HOME/.stew/bin to PATH"
fi

echo "\033[32m\033[1mSuccess:\033[0m Start a new terminal session to start using stew"
