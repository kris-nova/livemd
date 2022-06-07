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
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kris-nova/livemd/pkg/hackmd"

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
	title       string
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
				Name:      "notify",
				Aliases:   []string{"notification", "n"},
				Usage:     "Send notifications to configured backends.",
				UsageText: "live notify",
				Action: func(c *cli.Context) error {

					// Load notification system
					// Each notifier
					//  - Do Notify()     // Send the notification
					//  - Render Markdown // Return the pointer to the notification
					//
					// Note: The notifiers should just return [link](https://...) strings
					// Note: We have only ever used links in live streams in the past
					// Note: We can easily have a link system in *LiveMD

					return nil
				},
			},
			{
				//
				// [live hackmd]
				//
				Name:        "hackmd",
				Aliases:     []string{"h"},
				Usage:       "Manage local file sync with hackmd.",
				UsageText:   "live hackmd [commmand] [options]",
				Description: "Use this command to manage local state files with hackmd.",
				Subcommands: []*cli.Command{
					{
						Name:        "push",
						Usage:       "Push to remote.",
						UsageText:   "live hackmd push [options]",
						Description: "Use this to overwrite remote with local.",
						Flags:       GlobalFlags(StreamFlags([]cli.Flag{})),
						Action: func(c *cli.Context) error {
							// Check if state file exists
							if !livemd.FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] can not be found", cfg.filename)
							}
							l, err := livemd.Load(cfg.filename)
							if err != nil {
								return fmt.Errorf("file [%s] can not be loaded: %v", cfg.filename, err)
							}

							// Client
							if vars.hackmdToken == "" {
								return fmt.Errorf("empty token")
							}
							if vars.hackmdID == "" {
								return fmt.Errorf("empty hackmd id")
							}
							client := hackmd.New(vars.hackmdToken)
							_, err = client.GetNote(vars.hackmdID)
							if err != nil {
								return fmt.Errorf("unable to find note: %s: %v", vars.hackmdID, err)
							}
							note, err := livemd.ToNote(l, vars.hackmdID)
							if err != nil {
								return fmt.Errorf("unable to translate to note: %v", err)
							}
							_, err = client.UpdateNote(note)
							if err != nil {
								return err
							}
							logrus.Infof("Updated remote: https://hackmd.io/%s", vars.hackmdID)
							return nil
						},
					},
					{
						Name:        "pull",
						Usage:       "Pull from remote.",
						UsageText:   "live hackmd pull [options]",
						Description: "Use this to overwrite local with remote.",
						Flags:       GlobalFlags(StreamFlags([]cli.Flag{})),
						Action: func(c *cli.Context) error {
							// Client
							if vars.hackmdToken == "" {
								return fmt.Errorf("empty token")
							}
							if vars.hackmdID == "" {
								return fmt.Errorf("empty hackmd id")
							}
							client := hackmd.New(vars.hackmdToken)
							note, err := client.GetNote(vars.hackmdID)
							if err != nil {
								return fmt.Errorf("unable to find note: %s: %v", vars.hackmdID, err)
							}
							err = ioutil.WriteFile(cfg.filename, []byte(note.Content), livemd.DefaultMode)
							if err != nil {
								return err
							}
							logrus.Infof("Updated local: %s", cfg.filename)
							return nil
						},
					},
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
						UsageText:   "live stream new [options] <twitch-id>",
						Description: "Use this to create a new state file.",
						Flags:       GlobalFlags(StreamFlags([]cli.Flag{})),
						Action: func(c *cli.Context) error {
							twitchID := c.Args().Get(0)
							if twitchID == "" {
								logrus.Errorf("Missing <twitch-id>.")
								cli.ShowSubcommandHelp(c)
								return nil
							}

							// Check if state file exists
							if livemd.FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] exists", cfg.filename)
							}
							logrus.Infof("Creating new state: %s", cfg.filename)
							l := livemd.New(cfg.filename, vars.twitchChannel, twitchID)
							if strm.notify != "" {
								logrus.Infof("Setting notification string: %s", strm.notify)
								l.Notify = strm.notify
							}
							if strm.description != "" {
								logrus.Infof("Setting description string: %s", strm.notify)
								l.Notify = strm.description
							}

							logrus.Infof("Twitch Edit Dashboard: ")
							logrus.Infof(l.TwitchEditPage())

							logrus.Infof("Twitch Vidoe Page: ")
							logrus.Infof(l.TwitchVideoPage())

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

							// Todo use git to reconcile or something, IDK
							logrus.Warningf("Update will lose any remote information that is not templated!")

							// Check if state file exists
							if !livemd.FileExists(cfg.filename) {
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
								logrus.Infof("Setting description string: %s", strm.description)
								l.Description = strm.description
							}
							if strm.title != "" {
								logrus.Infof("Setting title string: %s", strm.title)
								l.Title = strm.title
							}
							logrus.Infof("Updating markdown. Rendering.")
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
							if !livemd.FileExists(cfg.filename) {
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
							return livemd.MoveFile(cfg.filename, writeFile)
						},
					},
					{
						Name:        "link",
						Usage:       "Add a link to the local state file.",
						UsageText:   "live stream link <title> <url>",
						Description: "Use this to add a link to the local state file.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {

							title := c.Args().Get(0)
							if title == "" {
								logrus.Errorf("Missing <title>.")
								cli.ShowSubcommandHelp(c)
								return nil
							}
							rawURL := c.Args().Get(1)
							if rawURL == "" {
								logrus.Errorf("Missing <url>.")
								cli.ShowSubcommandHelp(c)
								return nil
							}

							// Check if state file exists
							if !livemd.FileExists(cfg.filename) {
								return fmt.Errorf("file [%s] can not be found", cfg.filename)
							}
							l, err := livemd.Load(cfg.filename)
							if err != nil {
								return fmt.Errorf("file [%s] can not be loaded: %v", cfg.filename, err)
							}
							err = l.AddLink(title, rawURL)
							if err != nil {
								return fmt.Errorf("unable to add link: %v", err)
							}
							logrus.Infof("Saving: %s", cfg.filename)
							return l.Write()
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
			Destination: &strm.description,
			Value:       "",
		},
		&cli.StringFlag{
			Name:        "title",
			Aliases:     []string{"t"},
			Destination: &strm.title,
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
