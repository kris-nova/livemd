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
	"context"
	"fmt"

	"github.com/michimani/gotwi/tweet/managetweet/types"

	"github.com/bwmarrin/discordgo"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/sirupsen/logrus"
)

const (
	// MessageSizeLimit will be the maximum size of a message
	// that we can send with our notification system.
	// The limit is caused by Twitter's maximum tweet size.
	//
	// We believe Twitter has the SMALLEST constraint of all of our
	// integrations, and therefore sets the bar.
	MessageSizeLimit int = 280

	// TwitterName is the name to use for generating the URL in the logs
	TwitterName = "krisnova"
)

type Notifier struct {
	Message        string
	Discord        *discordgo.Session
	Twitter        *gotwi.Client
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
//
// The bot must be linked to a specific server which can be done via:
//   https://discord.com/oauth2/authorize?client_id=YOUR_ID_HERE&scope=bot&permissions=2048
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

func (n *Notifier) EnableTwitter(clientID, clientSecret string) error {
	twitter, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           clientID,
		OAuthTokenSecret:     clientSecret,
	})
	if err != nil {
		return fmt.Errorf("unable to enable twitter integration: %v", err)
	}
	n.Twitter = twitter
	if clientID == "" {
		return fmt.Errorf("empty twitter clientID")
	}
	if clientSecret == "" {
		return fmt.Errorf("empty twitter clientSecret")
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
		_, err := n.Discord.ChannelMessageSend(n.discordChannel, n.Message)
		if err != nil {
			logrus.Warningf("Discord notification failure: %v", err)
		} else {
			logrus.Info("Discord: Sent!")
		}
	}

	// Twitter
	if n.Twitter != nil {
		logrus.Info("Twitter: Dispatching...")
		res, err := managetweet.Create(context.TODO(), n.Twitter, &types.CreateInput{
			Text: gotwi.String(n.Message),
		})
		if err != nil {
			logrus.Warningf("Twitter notification failure: %v", err)
		} else {
			logrus.Info("Twitter: Sent!")
			logrus.Info("https://twitter.com/%s/status/", TwitterName, res.Data.ID)
		}
	}

	return nil

}
