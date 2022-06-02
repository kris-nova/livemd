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

package notify

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

const (
	// MessageSizeLimit will be the maximum size of a message
	// that we can send with our notification system.
	// The limit is caused by Twitter's maximum tweet size.
	//
	// We believe Twitter has the SMALLEST constraint of all of our
	// integrations, and therefore sets the bar.
	MessageSizeLimit int = 280
)

type Notifier struct {
	Message        string
	Discord        *discordgo.Session
	discordChannel string
}

func New(message string) *Notifier {
	return &Notifier{
		Message: message,
	}
}

// EnableDiscord will use the Bot API to communicate with Discord.
// You must pass the RAW token data to this function in order to create a new client.
//
// More: https://discord.com/developers/
func (n *Notifier) EnableDiscord(token string, channel string) error {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("unable to enable discord intergration: %v", err)
	}
	n.Discord = session
	n.discordChannel = channel
	if token == "" {
		return fmt.Errorf("empty discord token")
	}
	if channel == "" {
		return fmt.Errorf("empty channel ID")
	}
	return nil
}

func (n *Notifier) Notify() error {
	// Validation
	if len(n.Message) >= MessageSizeLimit {
		return fmt.Errorf("message size limit exceeded: %d: %d", len(n.Message), MessageSizeLimit)
	}

	// Discord
	if n.Discord != nil {
		logrus.Info("Discord: Dispatching...")
		n.Discord.Identify.Intents = 2048
		_, err := n.Discord.ChannelMessageSend(n.discordChannel, n.Message)
		if err != nil {
			logrus.Warningf("Discord notification failure: %v", err)
		}
	}

	return nil

}
