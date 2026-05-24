package cmd

import (
	"fmt"
	"os"
)

// GetZshCompletion returns the Zsh completion script as a string
func GetZshCompletion() string {
	return `#compdef stew
_stew_installed_binaries() {
    local -a binaries
    binaries=($(stew list 2>/dev/null | awk -F: 'NF {print $1}'))
    _describe 'installed binaries' binaries
}

_stew() {
    local context state state_descr line
    _arguments -C \
        '1: :->command' \
        '*: :->args'

    case $state in
        command)
            local -a commands
            commands=(
                'install:Install a binary from GitHub or URL'
                'i:Install a binary from GitHub or URL'
                'search:Search GitHub repos'
                's:Search GitHub repos'
                'browse:Browse releases and assets'
                'b:Browse releases and assets'
                'upgrade:Upgrade an installed binary'
                'up:Upgrade an installed binary'
                'uninstall:Uninstall a binary'
                'un:Uninstall a binary'
                'rename:Rename an installed binary'
                're:Rename an installed binary'
                'list:List installed binaries'
                'ls:List installed binaries'
                'config:Configure stew'
                'completion:Generate shell completion scripts'
                'help:Show help'
                'h:Show help'
            )
            _describe -t commands 'stew command' commands
            ;;
        args)
            case $line[1] in
                upgrade|up|uninstall|un|rename|re)
                    _stew_installed_binaries
                    ;;
            esac
            ;;
    esac
}
compdef _stew stew`
}

// RunCompletion prints the completion script to stdout for supported shells.
func RunCompletion(shell string) error {
	switch shell {
	case "zsh":
		fmt.Fprintln(os.Stdout, GetZshCompletion())
	default:
		return fmt.Errorf("unsupported shell: %s (only zsh is supported currently)", shell)
	}
	return nil
}
