# Configuration
`stew` can be configured with a `stew.config.json` file. The location of this file will depend on your OS:
|Linux/macOS | Windows |
| ------------ | ---------- |
| `$XDG_CONFIG_HOME/stew` or `~/.config/stew` | `~/AppData/Local/stew/Config` |

You can configure 2 aspects of `stew`:
1. The `stewPath`: this is where `stew` data is stored.
2. The `stewBinPath`: this is where `stew` installs binaries.
3. `excludeFromUpgradeAll`: this is the list of binaries that you don't want to keep upgraded during `stew upgrade --all`, perhaps if they have their own built in upgrade feature, or you want to pin a specific version.

The default locations for the first two are:
|                    | Linux/macOS | Windows |
| ------------ | ------------ | ---------- |
| `stewPath` | `$XDG_DATA_HOME/stew` or `~/.local/share/stew` | `~/AppData/Local/stew` |
| `stewBinPath` | `~/.local/bin` | `~/AppData/Local/stew/bin` |

There are multiple ways to configure these:
* When you first run `stew`, it will look for a `stew.config.json` file. If it cannot find one, then you will be prompted to set the configuration values.
* After `stew` is installed, you can use the `stew config` command to set the configuration values.
* At any time, you can manually create or edit the `stew.config.json` file. It should have values for `stewPath`, `stewBinPath`, and `excludeFromUpgradeAll`. 
