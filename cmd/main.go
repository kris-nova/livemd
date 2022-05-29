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
	"os"
	"time"

	"github.com/kris-nova/live"
	"github.com/kris-nova/live/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var cfg = &AppOptions{}

type AppOptions struct {
	verbose bool
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
		Name:     live.Name,
		Version:  live.Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  live.AuthorName,
				Email: live.AuthorEmail,
			},
		},
		Copyright: live.Copyright,
		HelpName:  live.Copyright,
		Usage:     "Collaborative Live Stream CLI Tool",
		UsageText: `live <cmd> <options> 
Use this program to perform tasks with Twitch, Hackmd, and YouTube.`,
		Commands: []*cli.Command{
			&cli.Command{},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Destination: &cfg.verbose,
			},
		},
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,
		Action: func(c *cli.Context) error {

			//
			novaObject := service.NewNova()
			return novaObject.Run()
			//

		},
	}
	Preloader()
	app.Run(os.Args)
}

// Preloader will run for ALL commands, and is used
// to initalize the runtime environments of the program.
func Preloader() {
	/* Flag parsing */
	if cfg.verbose {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
