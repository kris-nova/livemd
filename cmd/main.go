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

	"github.com/kris-nova/livemd"

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

func main() {
	/* Change version to -V */
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "The version of the program.",
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s
{{.Usage}}

Commands{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}

Options
   {{range .VisibleFlags}}{{.}}
   {{end}}
`, livemd.Banner())

	app := &cli.App{
		Name:     "live",
		Version:  livemd.Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  livemd.AuthorName,
				Email: livemd.AuthorEmail,
			},
		},
		Copyright: livemd.Copyright,
		HelpName:  livemd.Name,
		Usage:     "Centralized live streaming meta data in markdown.",
		Commands: []*cli.Command{
			{
				Name:      "notification",
				Aliases:   []string{"notif", "n", "notify"},
				Usage:     "Send notifications to configured backends.",
				UsageText: "live notify <message>",
				Action: func(c *cli.Context) error {
					//message := c.Args().Get(0)
					//if message == "" {
					//	fmt.Println(livemd.Banner())
					//	cli.ShowSubcommandHelp(c)
					//	return nil
					//}
					//notifier := notify.New(message)
					//logrus.Infof("=== Starting Notification Bus ===")
					//var err error
					//if cfg.discordToken != "" {
					//	err = notifier.EnableDiscord(cfg.discordToken, cfg.discordChannel)
					//	if err != nil {
					//		return fmt.Errorf("failed enabling discord: %v", err)
					//	}
					//}
					//if cfg.twitterApiKeySecret != "" {
					//	err = notifier.EnableTwitter(cfg.twitterAccessToken, cfg.twitterAccessTokenSecret, cfg.twitterApiKey, cfg.twitterApiKeySecret)
					//	if err != nil {
					//		return fmt.Errorf("failed enabling twitter: %v", err)
					//	}
					//}
					//
					//if cfg.mastodonUsername != "" {
					//	err = notifier.EnableMastodon(cfg.mastodonServer, cfg.mastodonAccessToken, cfg.mastodonClientID, cfg.mastodonClientSecret, cfg.mastodonUsername, cfg.mastodonPassword)
					//	if err != nil {
					//		return fmt.Errorf("failed enabling mastodon: %v", err)
					//	}
					//}
					//
					//// Run the notifications system
					//err = notifier.Notify()
					//logrus.Infof("=== Stopping Notification Bus ===")
					//if err != nil {
					//	return err
					//}

					return nil
				},
			},
			{
				Name:        "stream",
				Aliases:     []string{"s"},
				Usage:       "Manage local live stream records.",
				UsageText:   "live stream <title>",
				Description: "Use this command to manage local records.",
				Subcommands: []*cli.Command{
					{
						Name:        "push",
						Usage:       "Push up the local livemd.md",
						UsageText:   "live stream <title>",
						Description: "Use this sync from remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error { //
							return nil
						},
					},
					{
						Name:        "pull",
						Usage:       "Pull down to the local livemd.md",
						UsageText:   "live stream <title>",
						Description: "Use this command to overwrite remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {
							return nil
						},
					},
				},
				Flags: GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					title := c.Args().Get(0)
					if title == "" {
						cli.ShowAppHelp(c)
						return nil
					}

					return nil
				},
			},
		},
		Flags: GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
	}
	var err error
	err = LoadEnvironmentalVariables()
	if err != nil {
		logrus.Errorf("Failed Loading Environment: %v", err)
		os.Exit(80)
	}
	if cfg.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	err = app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Runtime: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
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
			Value:       DefaultFile,
		},
	}
	for _, gf := range g {
		c = append(c, gf)
	}
	return c
}
