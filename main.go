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
		Version: "v0.1.0",
		Commands: []cli.Command{
			{
				Name:    "install",
				Usage:   "Install a binary",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					cmd.Install(c.Args())
					return nil
				},
			},
			{
				Name:    "search",
				Usage:   "Search for a project on GitHub",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					cmd.Search(c.Args().First())
					return nil
				},
			},
			{
				Name:    "browse",
				Usage:   "Browse releases and assets in a GitHub repo",
				Aliases: []string{"b"},
				Action: func(c *cli.Context) error {
					cmd.Browse(c.Args().First())
					return nil
				},
			},
			{
				Name:    "upgrade",
				Usage:   "Upgrade a binary to the latest available in the GitHub repo",
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
				Usage:   "Uninstall a binary",
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
				Name:    "list",
				Usage:   "List installed binaries",
				Aliases: []string{"ls"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "tags",
						Usage: "include the version tags",
					},
					&cli.BoolFlag{
						Name:  "assets",
						Usage: "include the assets and version tags",
					},
				},
				Action: func(c *cli.Context) error {
					cmd.List(c.Bool("tags"), c.Bool("assets"))
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
