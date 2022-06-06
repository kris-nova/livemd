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
	"path/filepath"
	"time"

	"github.com/kris-nova/livemd"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	DefaultFile   string = "live.md"
	DefaultOutput string = "archive"
)

var cfg = &AppOptions{}

type AppOptions struct {
	verbose  bool
	filename string
	output   string
}

var strm = &StreamOptions{}

type StreamOptions struct {
	notify      string
	description string
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

	cli.SubcommandHelpTemplate = fmt.Sprintf(`%s
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
				//
				// [live stream]
				//
				// This subcommand will only mutate local text.
				// By design this subcommand will NEVER reach out
				// to interface with other APIs.
				//
				Name:        "stream",
				Aliases:     []string{"s"},
				Usage:       "Manage local live stream state files.",
				UsageText:   "live stream <title>",
				Description: "Use this command to manage local state files.",
				Subcommands: []*cli.Command{
					{
						Name:        "new",
						Usage:       "Create a new local state file.",
						UsageText:   "live stream new [options] <title>",
						Description: "Use this to create a new state file.",
						Flags:       GlobalFlags(StreamFlags([]cli.Flag{})),
						Action: func(c *cli.Context) error {
							title := c.Args().Get(0)
							if title == "" {
								logrus.Errorf("Missing <title>.")
								cli.ShowSubcommandHelp(c)
								return nil
							}

							// Check if state file exists
							if FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] exists", cfg.filename)
							}
							logrus.Infof("Creating new state: %s", cfg.filename)
							l := livemd.New(cfg.filename)
							l.Title = title
							if strm.notify != "" {
								logrus.Infof("Setting notification string: %s", strm.notify)
								l.Notify = strm.notify
							}
							if strm.description != "" {
								logrus.Infof("Setting description string: %s", strm.notify)
								l.Notify = strm.description
							}
							return l.Write()
						},
					},
					{
						Name:        "update",
						Usage:       "Update fields in a local state file",
						UsageText:   "live stream update [options]",
						Description: "Use this to update an existing state file.",
						Flags:       GlobalFlags(StreamFlags([]cli.Flag{})),
						Action: func(c *cli.Context) error {

							// Check if state file exists
							if !FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] can not be found", cfg.filename)
							}
							logrus.Infof("Updating state: %s", cfg.filename)
							l, err := livemd.Load(cfg.filename)
							if err != nil {
								return fmt.Errorf("file [%s] can not be loaded: %v", cfg.filename, err)
							}
							if strm.notify != "" {
								logrus.Infof("Setting notification string: %s", strm.notify)
								l.Notify = strm.notify
							}
							if strm.description != "" {
								logrus.Infof("Setting description string: %s", strm.notify)
								l.Notify = strm.description
							}
							return l.Write()
						},
					},
					{
						Name:        "archive",
						Usage:       "Archive a local state file",
						UsageText:   "live stream archive [options]",
						Description: "Use this to update archive a local state file.",
						Flags: GlobalFlags([]cli.Flag{
							&cli.StringFlag{
								Name:        "output",
								Aliases:     []string{"o"},
								Destination: &cfg.output,
								Value:       DefaultOutput,
							},
						}),
						Action: func(c *cli.Context) error {

							// Check if state file exists
							if !FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] can not be found", cfg.filename)
							}
							l, err := livemd.Load(cfg.filename)
							if err != nil {
								return fmt.Errorf("file [%s] can not be loaded: %v", cfg.filename, err)
							}
							titleName := l.CoalesceDateName()
							writeFile := filepath.Join(cfg.output, titleName)
							logrus.Infof("Archiving state: %s", writeFile)

							// By design this should NEVER l.Write()
							// We can use the variables defined in l
							// However we should NEVER mutate the source during an archive!
							return MoveFile(cfg.filename, writeFile)
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

func StreamFlags(c []cli.Flag) []cli.Flag {
	g := []cli.Flag{
		&cli.StringFlag{
			Name:        "notify",
			Aliases:     []string{"n"},
			Destination: &strm.notify,
			Value:       "",
		},
		&cli.StringFlag{
			Name:        "description",
			Aliases:     []string{"d"},
			Destination: &strm.notify,
			Value:       "",
		},
	}
	for _, gf := range g {
		c = append(c, gf)
	}
	return c
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

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// MoveFile is primarily used for Archiving.
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}
