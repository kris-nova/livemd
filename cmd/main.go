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
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"time"

	"github.com/kris-nova/live/pkg/notify"

	"github.com/kris-nova/live/pkg/discord"

	"github.com/kris-nova/live/pkg/embedmd"

	"github.com/kris-nova/live/pkg/hackmd"

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
	verbose        bool
	filename       string
	hackmdToken    string
	hackmdID       string
	discordToken   string
	discordChannel string
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
				Name:      "notification",
				Aliases:   []string{"notif", "n", "notify"},
				Usage:     "Send notifications to configured backends.",
				UsageText: "live notify <message>",
				Action: func(c *cli.Context) error {
					message := c.Args().Get(0)
					if message == "" {
						fmt.Println(live.Banner())
						cli.ShowSubcommandHelp(c)
						return nil
					}
					notifier := notify.New(message)
					var err error
					err = notifier.EnableDiscord(cfg.discordToken, cfg.discordChannel)
					if err != nil {
						return fmt.Errorf("invalid discord token: %v", err)
					} else {
						logrus.Infof("Discord: Enabled")
					}

					// Run the notifications system
					return notifier.Notify()
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
						Usage:       "Push up the local live.md",
						UsageText:   "live stream <title>",
						Description: "Use this sync from remote.",
						Flags:       GlobalFlags([]cli.Flag{}),
						Action: func(c *cli.Context) error {
							x, err := livemd.FromLocal(cfg.filename)
							if err != nil {
								return fmt.Errorf("unable to find local: %s: %v", cfg.filename, err)
							}
							client := hackmd.New(cfg.hackmdToken)
							if cfg.hackmdID == "" {
								return fmt.Errorf("empty hackmd id")
							}
							err = x.Sync(cfg.filename)
							if err != nil {
								return fmt.Errorf("sync: %v", err)
							}
							_, err = client.GetNote(cfg.hackmdID)
							if err != nil {
								// Does not exist
								return fmt.Errorf("unable to find note: %s: %v", cfg.hackmdID, err)
							}
							_, err = client.UpdateNote(x.ToHackMD(cfg.hackmdID))
							if err != nil {
								return fmt.Errorf("unable to push: %v", err)
							}
							logrus.Infof("Saved hackmd   : %s", cfg.hackmdID)
							logrus.Infof("Saved filename : %s", cfg.filename)
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
							client := hackmd.New(cfg.hackmdToken)
							if cfg.hackmdID == "" {
								return fmt.Errorf("empty HACKMD_ID")
							}
							y, err := client.GetNote(cfg.hackmdID)
							if err != nil {
								return fmt.Errorf("unable to find local: %s: %v", cfg.filename, err)
							}
							x := &livemd.LiveMD{}
							err = embedmd.Unmarshal([]byte(y.Content), x)
							if err != nil {
								return fmt.Errorf("invalid remote, unable to Unmarshal: %v", err)
							}
							x.SyncRaw([]byte(y.Content))
							err = x.Write(cfg.filename)
							if err != nil {
								return fmt.Errorf("unable to write: %v", err)
							}
							logrus.Infof("Saved hackmd   : %s", cfg.hackmdID)
							logrus.Infof("Saved filename : %s", cfg.filename)
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
							x, err := livemd.FromLocal(DefaultFile)
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
					x := livemd.New(title)
					x.HackMDID = cfg.hackmdID
					err := x.Write(cfg.filename)
					if err != nil {
						return fmt.Errorf("unable to write local: %v", err)
					}
					client := hackmd.New(cfg.hackmdToken)
					note := x.ToHackMD(cfg.hackmdID)
					if note.ID == "" {
						logrus.Infof("Creating new hackMD note: %s", note.Title)
						note, err = client.CreateNote(note)
						x.HackMDID = note.ID
						defer x.Write(cfg.filename)
					} else {
						logrus.Infof("Updating hackMD note: %s", note.Title)
						note, err = client.UpdateNote(note)
					}
					if err != nil {
						return fmt.Errorf("unable to save to hackmd: %v", err)
					}
					logrus.Infof("Successful save: %s", note.ID)
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
	err := Validation()
	if err != nil {
		logrus.Errorf("Failed Validation: %v", err)
		os.Exit(99)
	}
	err = app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Runtime: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func Validation() error {
	cfg.hackmdToken = os.Getenv(hackmd.EnvironmentalVariableHackMDToken)
	if cfg.hackmdToken == "" {
		return fmt.Errorf("empty environmental variable [%s]", hackmd.EnvironmentalVariableHackMDToken)
	}
	return nil
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

	// ---
	// Environmental Variable Parsing
	// ---

	cfg.hackmdID = os.Getenv(hackmd.EnvironmentalVariableHackMDID)
	if cfg.hackmdID != "" {
		logrus.Infof("Loading HackMD ID: %s", cfg.hackmdID)
	}

	cfg.discordToken = os.Getenv(discord.EnvironmentalVariableDiscordToken)
	if cfg.discordToken != "" {
		logrus.Infof("Loading Discord Token: **********")
	}

	cfg.discordChannel = os.Getenv(discord.EnvironmentalVariableDiscordChannel)
	if cfg.discordToken != "" {
		logrus.Infof("Loading Discord Channel: %s", cfg.discordChannel)
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
			Value:       DefaultFile,
		},
	}
	for _, gf := range g {
		c = append(c, gf)
	}
	return c
}

func Compare(a, b *hackmd.Note) bool {

	h1 := hash(string(a.Content))
	h2 := hash(string(b.Content))
	logrus.Infof("Compare %d:%d", h1, h2)
	return h1 == h2
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Pause() {
	fmt.Print("Press any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
