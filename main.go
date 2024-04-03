package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v3"

	"github.com/marwanhawari/stew/cmd"
	stew "github.com/marwanhawari/stew/lib"
)

func main() {
	app := &cli.Command{
		Name:                  "stew",
		EnableShellCompletion: true,
		Version:               "v0.3.0",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Usage:   "Install a binary. The input can be a GitHub repo or a URL. [Ex: stew install marwanhawari/ppath]",
				Aliases: []string{"i"},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Install(c.Args().Slice())
					return nil
				},
			},
			{
				Name:    "search",
				Usage:   "Search for a GitHub repo then browse the selected repo's releases and assets. [Ex: stew search ripgrep]",
				Aliases: []string{"s"},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Search(c.Args().First())
					return nil
				},
			},
			{
				Name:    "browse",
				Usage:   "Browse the releases and assets from a GitHub repo. [Ex: stew browse marwanhawari/ppath]",
				Aliases: []string{"b"},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Browse(c.Args().First())
					return nil
				},
			},
			{
				Name:    "upgrade",
				Usage:   "Upgrade a binary to the latest available in the GitHub repo. Use the name of the installed binary. [Ex: stew upgrade fzf]",
				Aliases: []string{"up"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Upgrade all binaries",
					},
				},
				ShellComplete: listInstalledBinaries,
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Upgrade(c.Bool("all"), c.Args().First())
					return nil
				},
			},
			{
				Name:    "uninstall",
				Usage:   "Uninstall a binary. Use the name of the installed binary. [Ex: stew uninstall fzf]",
				Aliases: []string{"un"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Uninstall all binaries",
					},
				},
				ShellComplete: listInstalledBinaries,
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Uninstall(c.Bool("all"), c.Args().First())
					return nil
				},
			},
			{
				Name:          "rename",
				Usage:         "Rename an installed binary using an interactive UI. [Ex: stew rename fzf]",
				Aliases:       []string{"re"},
				ShellComplete: listInstalledBinaries,
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Rename(c.Args().First())
					return nil
				},
			},
			{
				Name:    "list",
				Usage:   "List installed binaries [Ex: stew list]",
				Aliases: []string{"ls"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "tags",
						Usage: "include the version tags",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.List(c.Bool("tags"))
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "Configure the stew file paths using an interactive UI. [Ex: stew config]",
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.Config()
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func listInstalledBinaries(ctx context.Context, cmd *cli.Command) {
	configPath, err := stew.GetStewConfigFilePath(runtime.GOOS)
	if err != nil {
		return
	}
	configExists, err := stew.PathExists(configPath)
	if err != nil || !configExists {
		return
	}
	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	stewLockFilePath := systemInfo.StewLockFilePath

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	if err != nil {
		return // no lockfile
	}
	for _, pkg := range lockFile.Packages {
		fmt.Println(pkg.Binary)
	}
}
