package main

import (
	"log"
	"os"

	"github.com/marwanhawari/stew/cmd"
	"github.com/urfave/cli"
)

func main() {

	app := &cli.App{
		Name:    "stew",
		Version: "v0.5.0",
		Commands: []cli.Command{
			{
				Name:    "install",
				Usage:   "Install a binary. The input can be a GitHub repo or a URL. [Ex: stew install marwanhawari/ppath]",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					cmd.Install(c.Args().First())
					return nil
				},
			},
			{
				Name:    "search",
				Usage:   "Search for a GitHub repo then browse the selected repo's releases and assets. [Ex: stew search ripgrep]",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					cmd.Search(c.Args())
					return nil
				},
			},
			{
				Name:    "browse",
				Usage:   "Browse the releases and assets from a GitHub repo. [Ex: stew browse marwanhawari/ppath]",
				Aliases: []string{"b"},
				Action: func(c *cli.Context) error {
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
				Action: func(c *cli.Context) error {
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
				Action: func(c *cli.Context) error {
					cmd.Uninstall(c.Bool("all"), c.Args().First())
					return nil
				},
			},
			{
				Name:    "rename",
				Usage:   "Rename an installed binary using an interactive UI. [Ex: stew rename fzf]",
				Aliases: []string{"re"},
				Action: func(c *cli.Context) error {
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
				Action: func(c *cli.Context) error {
					cmd.List(c.Bool("tags"))
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "Configure stew using an interactive UI. [Ex: stew config]",
				Action: func(c *cli.Context) error {
					cmd.Config()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
