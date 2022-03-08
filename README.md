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
* Install binaries from GitHub releases or directly from URLs
* Easily distribute binaries across teams and private repositories.
* Get the latest releases ahead of other package managers.
* Rapidly browse, install, and experiment with different projects.
* [Configure](https://github.com/marwanhawari/stew/blob/main/config.md) where to install binaries.
* No need for `sudo`.
* Portable [`Stewfile`](https://github.com/marwanhawari/stew/blob/main/examples/Stewfile) with optional pinned versioning.

![demo](https://github.com/marwanhawari/stew/raw/main/assets/demo.gif)

# Installation
Stew supports Linux, macOS, and Windows.

### Install a pre-compiled binary:

* Install using `curl`:
```
curl -fsSL https://raw.githubusercontent.com/marwanhawari/stew/main/install.sh | bash
```

* Install using `brew`:
```
brew install marwanhawari/tap/stew
```

* Download a pre-compiled binary from the [releases page](https://github.com/marwanhawari/stew/releases).

### Install from source:

* Install the latest released version:
```
go install github.com/marwanhawari/stew@latest
```

* Install the latest unreleased source:
```
git clone https://github.com/marwanhawari/stew
cd stew
go install .
```

# Usage
### Install
```sh
# Install from GitHub releases
stew install junegunn/fzf              # Install the latest release
stew install junegunn/fzf@0.27.1       # Install a specific, tagged version
stew install junefunn/fzf sharkdp/fd   # Install multiple binaries in a single command

# Install directly from a URL
stew install https://github.com/cli/cli/releases/download/v2.4.0/gh_2.4.0_macOS_amd64.tar.gz

# Install from an Stewfile
stew install Stewfile
```

### Search
```sh
# Search for a GitHub repo and browse its contents with a terminal UI
stew search ripgrep
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
stew list --tags --assets > Stewfile   # Pin tags and assets
```

### Config
```sh
# Configure the stew file paths using an interactive UI
stew config           # Automatically updates the stew.config.json
```

# FAQ
### Why couldn't `stew` automatically find any binaries for X repo?
The repo probably uses an unconventional naming scheme for their binaries. You can always manually select the release asset.

### Will `stew` work with private GitHub repositories?
Yes, `stew` will automatically detect if you have a `GITHUB_TOKEN` environment variable and allow you to access binaries from your private repositories.

### Where does `stew` install binaries?
The default installation path will depend on your OS:
| Linux/macOS | Windows |
| ------------ | ---------- |
| `~/.local/bin` | `~/AppData/Local/stew/bin` |

However, this location can be [configured](https://github.com/marwanhawari/stew/blob/main/config.md).

Make sure that the installation path is in your `PATH` environment variable. Otherwise, you won't be able to use any of the binaries installed by `stew`.

### I've installed `stew` but the command is still not found.
If you install `stew` with the [curl install script](https://github.com/marwanhawari/stew/blob/main/install.sh), the `stew` binary will be installed in the `stewBinPath` you provided. The install script also attempts to add the `stewBinPath` to `PATH` in your `.zshrc` or `.bashrc` file. You will need to start a new terminal session for the changes to take effect. Make sure that the `stewBinPath` is in your `PATH` environment variable.
