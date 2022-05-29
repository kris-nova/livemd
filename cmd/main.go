/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kris-nova/live/pkg/livemd"

	"github.com/kris-nova/live"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	DefaultFile string = "live.md"
)

var cfg = &AppOptions{}

type AppOptions struct {
	verbose  bool
	filename string
}

// # Edit ./live.md
// live stream <title>    # Create a new live stream (hackmd)
// live stream push       # Sync local changes to hackmd
// live stream pull       # Overwrite local changes to hackmd
//
// # Update firebot with new hackmd URL (TODO Automate)
// live twitch push       # Sync local file to twitch API
// live twitch pull       # Overwrite local changes from twitch API
// live twitch export     # Export twitch episode to YouTube

func main() {
	/* Change version to -V */
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "The version of the program.",
	}
	app := &cli.App{
		Name:     "live",
		Version:  live.Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  live.AuthorName,
				Email: live.AuthorEmail,
			},
		},
		Copyright: live.Copyright,
		HelpName:  live.Name,
		Usage:     "Collaborative live stream CLI tool writte by Kris Nóva.",
		UsageText: `live <cmd> <options> 
Use this program to perform tasks with Twitch, Hackmd, and YouTube.`,
		Commands: []*cli.Command{
			{
				Name:        "stream",
				Aliases:     []string{"s"},
				Usage:       "Manage local live stream records.",
				UsageText:   "live stream <title>",
				Description: "Use this command to manage local records.",
				Subcommands: []*cli.Command{
					{
						Name:        "push",
						Usage:       "Push up the local live.md",
						UsageText:   "live stream <title>",
						Description: "Use this sync from remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {
							return nil
						},
					},
					{
						Name:        "pull",
						Usage:       "Pull down to the local live.md",
						UsageText:   "live stream <title>",
						Description: "Use this command to overwrite remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {
							return nil
						},
					},
					{
						Name:        "status",
						Usage:       "Show status of local file.",
						UsageText:   "live stream status",
						Description: "Use this command to overwrite remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {

							// Status always comes from local state
							x, err := livemd.FromFile(DefaultFile)
							if err != nil {
								return fmt.Errorf("unable to open %s: %v", DefaultFile, err)
							}
							status := x.Status()
							fmt.Print(status)
							return nil
						},
					},
				},
				Flags: GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					title := c.Args().Get(0)
					if title == "" {
						fmt.Println(live.Banner())
						cli.ShowSubcommandHelp(c)
						return nil
					}

					// New with name
					logrus.Infof("Creating New Stream \"%s\"", title)
					return nil
				},
			},
		},
		Flags: GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			fmt.Println(live.Banner())
			cli.ShowSubcommandHelp(c)
			return nil
		},
	}
	Preloader()
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Runtime: %v", err)
	}
}

// Preloader will run for ALL commands, and is used
// to initalize the runtime environments of the program.
func Preloader() {
	/* Flag parsing */
	if cfg.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func GlobalFlags(c []cli.Flag) []cli.Flag {
	g := []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Destination: &cfg.verbose,
		},
		&cli.StringFlag{
			Name:        "filename",
			Aliases:     []string{"f"},
			Destination: &cfg.filename,
		},
	}
	for _, gf := range g {
		c = append(c, gf)
	}
	return c
}
