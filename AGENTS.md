# AGENTS.md

## Essential Commands

- **Build**: `go build -o stew main.go`
- **Test**: `go test -v ./...`
- **Run**: `./stew [command]`
- **Lint**: No dedicated linter configured in the project, but `go vet ./...` is recommended.
- **Coverage**: `go test -v ./... --coverprofile coverage.out && go tool cover --func coverage.out`
- **Completion**: `stew completion zsh` - Generates Zsh shell completion script.

## Code Organization

- `main.go`: Entry point, defines CLI structure using `github.com/urfave/cli`.
- `cmd/`: Command implementations (cli logic).
- `lib/`: Core library logic (business logic, GitHub API, file operations, configuration).
- `constants/`: Shared constants like regex patterns for OS/Arch detection.
- `pkg/commands/`: Appears mostly empty or reserved for further command abstractions.

## Architecture & Flow

`stew` is a binary manager that installs executables directly from GitHub releases or URLs.

1.  **Initialization**: `stew.Initialize()` (in `lib/config.go`) detects OS/Arch and loads config.
2.  **Configuration**: Stored in `~/.config/stew/stew.config.json` (on Linux). Manages `stewPath` and `stewBinPath`.
3.  **State Tracking**: Uses a `Stewfile.lock.json` in the `stewPath` to track installed binaries, their versions, and hashes.
4.  **Installation Flow**:
    - Parse input (owner/repo or URL).
    - Query GitHub API for releases/assets if applicable.
    - Match assets using regex patterns in `constants/` against system OS/Arch.
    - Download and extract (using `github.com/mholt/archiver`).
    - Detect the binary in the extracted files (looks for files with executable permissions).
    - Move binary to `stewBinPath`.
    - Update `Stewfile.lock.json`.

## Conventions & Patterns

- **Error Handling**: Custom error types are defined in `lib/errors.go`. Most CLI commands use `stew.CatchAndExit(err)` for consistent error reporting.
- **UI**: Interactive prompts are common, using `github.com/AlecAivazis/survey/v2` (helper functions in `lib/ui.go`).
- **Styles**: ANSI colors are applied via `github.com/gookit/color` (abstractions in `constants/constants.go`).
- **Naming**: Follows standard Go naming conventions. `lib` package (aliased as `stew` in some places) contains the bulk of the internal API.

## Testing Approach

- Tests are located alongside source files (e.g., `lib/util_test.go`).
- Uses standard `go test`.
- GitHub Actions workflow (`.github/workflows/test.yml`) runs tests on every push/PR.

## Important Gotchas

- **Self-Installation**: Stew explicitly prevents self-installation or self-upgrade (`lib/errors.go`, `SelfInstallError`).
- **Binary Detection**: Stew relies on finding exactly one executable file in the downloaded asset. If it finds multiple or none, it prompts the user.
- **OS/Arch Matching**: Done via regex. See `constants/constants.go` for the exact patterns used to match asset names.
- **Lockfile**: The `Stewfile.lock.json` is critical for `upgrade` and `uninstall` commands. If it gets out of sync with the actual `bin` directory, these commands might fail or behave unexpectedly.
- **Zsh Completions**: The `completion` command generates a Zsh script. It currently only supports Zsh. It relies on `stew list` output to populate installed binaries for command arguments (like `upgrade` or `uninstall`).
