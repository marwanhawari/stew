<p align="center">
  <img width=30% height=auto src="https://github.com/marwanhawari/stew/raw/main/assets/stew.png" alt="stew icon"/>
</p>

<h3 align="center">stew</h3>
<p align="center">
  An independent package manager for compiled binaries.
</p>
<p align="center">

  <a href="https://github.com/marwanhawari/stew/actions/">
    <img src="https://github.com/marwanhawari/stew/actions/workflows/test.yml/badge.svg" alt="build status"/>
  </a>
  
  <a href="https://goreportcard.com/report/github.com/marwanhawari/stew">
    <img src="https://goreportcard.com/badge/github.com/marwanhawari/stew" alt="go report card"/>
  </a>

  <a href='https://coveralls.io/github/marwanhawari/stew?branch=main'>
    <img src='https://coveralls.io/repos/github/marwanhawari/stew/badge.svg?branch=main' alt='Coverage Status'/>
  </a>

  <a href="https://pkg.go.dev/github.com/marwanhawari/stew">
    <img src="https://pkg.go.dev/badge/github.com/marwanhawari/stew.svg" alt="pkg.go.dev reference"/>
  </a>

  <a href="https://github.com/avelino/awesome-go">
    <img src="https://awesome.re/mentioned-badge.svg" alt="Mentioned in Awesome Go"/>
  </a>
  
</p>


# Features
* Install binaries from GitHub releases or directly from URLs.
* Easily distribute binaries across teams and private repositories.
* Get the latest releases ahead of other package managers.
* Rapidly browse, install, and experiment with different projects.
* [Configure](https://github.com/marwanhawari/stew/blob/main/config.md) where to install binaries.
* No need for `sudo`.
* Just a single binary with 0 dependencies.
* Portable [`Stewfile`](https://github.com/marwanhawari/stew/blob/main/examples/Stewfile) with optional pinned versioning.
* Headless batch installs from a `Stewfile.lock.json` file.

![demo](https://github.com/marwanhawari/stew/raw/main/assets/demo.gif)

# Installation
Stew supports macOS, Linux, and Windows.

### Install using a package manager
<details>
  <summary>macOS</summary>

  ```sh
  brew install marwanhawari/tap/stew
  ```
</details>

<details>
  <summary>Arch</summary>

  ```sh
  git clone https://aur.archlinux.org/stew.git
  cd stew
  makepkg -sric
  ```
</details>

### Download a compiled binary
Compiled binaries can be downloaded from the [releases page](https://github.com/marwanhawari/stew/releases).

### Install using Go
<details>
  <summary>Install the latest released version</summary>

  ```sh
  go install github.com/marwanhawari/stew@latest
  ```
</details>

<details>
  <summary>Install the latest unreleased source</summary>

  ```sh
  git clone https://github.com/marwanhawari/stew
  cd stew
  go install .
  ```
</details>

# Usage
### Install
```sh
# Install from GitHub releases
stew install junegunn/fzf              # Install the latest release
stew install junegunn/fzf@0.27.1       # Install a specific, tagged version

# Install directly from a URL
stew install https://github.com/cli/cli/releases/download/v2.4.0/gh_2.4.0_macOS_amd64.tar.gz

# Install from an Stewfile
stew install Stewfile

# Install headlessly from a Stewfile.lock.json
stew install Stewfile.lock.json

# Install mutliple binaries per repo/asset
stew install astral-sh/uv     # Install uv the first time
stew install astral-sh/uv     # Install uvx the second time
```

### Search
```sh
# Search for a GitHub repo and browse its contents with a terminal UI
stew search ripgrep
stew search fzf user:junegunn language:go    # Use GitHub search syntax
```

### Browse
```sh
# Browse a specific GitHub repo's releases and assets with a terminal UI
stew browse sharkdp/hyperfine
```

### Upgrade
```sh
# Upgrade a binary to its latest version. Not for binaries installed from a URL.
stew upgrade rg           # Upgrade using the name of the binary directly
stew upgrade --all        # Upgrade all binaries
```

### Uninstall
```sh
# Uninstall a binary
stew uninstall rg         # Uninstall using the name of the binary directly
stew uninstall --all      # Uninstall all binaries
```

### Rename
```sh
# Rename an installed binary using an interactive UI
stew rename rg            # Rename using the name of the binary directly
```

### List
```sh
# List installed binaries
stew list                              # Print to console
stew list > Stewfile                   # Create an Stewfile without pinned tags
stew list --tags > Stewfile            # Pin tags
```

### Config
```sh
# Configure the stew file paths using an interactive UI
stew config           # Automatically updates the stew.config.json
```

# Configuration

`stew` can be configured with a `stew.config.json` file. The location of this file will also depend on your OS:
|Linux/macOS | Windows |
| ------------ | ---------- |
| `$XDG_CONFIG_HOME/stew` or `~/.config/stew` | `~/AppData/Local/stew/Config` |

You can configure 2 aspects of `stew`:
1. The `stewPath`: this is where `stew` data is stored.
2. The `stewBinPath`: this is where `stew` installs binaries
3. `excludeFromUpgradeAll`: this is the list of binaries that you don't want to keep upgraded during `stew upgrade --all`, perhaps if they have their own built in upgrade feature, or you want to pin a specific version.

The default locations for the `stewPath` and `stewBinPath` are:
|                    | Linux/macOS | Windows |
| ------------ | ------------ | ---------- |
| `stewPath` | `$XDG_DATA_HOME/stew` or `~/.local/share/stew` | `~/AppData/Local/stew` |
| `stewBinPath` | `~/.local/bin` | `~/AppData/Local/stew/bin` |

There are multiple ways to configure these:
* When you first run `stew`, it will look for a `stew.config.json` file. If it cannot find one, then you will be prompted to set the configuration values.
* After `stew` is installed, you can use the `stew config` command to set the configuration values.
* At any time, you can manually create or edit the `stew.config.json` file. It should have values for `stewPath`, `stewBinPath`, and `excludeFromUpgradeAll`. 

Make sure that the installation path is in your `PATH` environment variable. Otherwise, you won't be able to use any of the binaries installed by `stew`.


# FAQ
### Why couldn't `stew` automatically find any binaries for X repo?
The repo probably uses an unconventional naming scheme for their binaries. You can always manually select the release asset.

### Will `stew` work with private GitHub repositories?
Yes, `stew` will automatically detect if you have a `GITHUB_TOKEN` environment variable and allow you to access binaries from your private repositories.
