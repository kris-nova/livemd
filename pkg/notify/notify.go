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
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/chimeracoder/anaconda"
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

	EnableDiscordSend bool = false
	EnableTwitterSend bool = false
)

type Notifier struct {
	Message        string
	Discord        *discordgo.Session
	discordChannel string
	Twitter        *anaconda.TwitterApi
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

func (n *Notifier) EnableTwitter(accessToken, accessTokenSecret, consumerKey, consumerKeySecret string) error {
	if accessToken == "" {
		return fmt.Errorf("empty accessToken")
	}
	if accessTokenSecret == "" {
		return fmt.Errorf("empty accessTokenSecret")
	}
	if consumerKey == "" {
		return fmt.Errorf("empty consumerKey")
	}
	if consumerKeySecret == "" {
		return fmt.Errorf("empty consumerKeySecret")
	}
	twitter := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerKeySecret)
	self, err := twitter.GetSelf(url.Values{})
	if err != nil {
		return fmt.Errorf("unable to enable twitter integration: %v", err)
	}
	n.Twitter = twitter
	logrus.Infof("Twitter authenticated! User: %s", self.Name)
	return nil
}

func (n *Notifier) Notify() error {
	// Validation
	if len(n.Message) >= MessageSizeLimit {
		return fmt.Errorf("message size limit exceeded: %d: %d", len(n.Message), MessageSizeLimit)
	}

	// Discord
	if n.Discord != nil && EnableDiscordSend {
		logrus.Info("Discord: Dispatching...")
		_, err := n.Discord.ChannelMessageSend(n.discordChannel, n.Message)
		if err != nil {
			logrus.Warningf("Discord notification failure: %v", err)
		} else {
			logrus.Info("Discord: Sent!")
		}
	}

	// Twitter
	if n.Twitter != nil && EnableTwitterSend {
		logrus.Info("Twitter: Dispatching...")
		tweet, err := n.Twitter.PostTweet(n.Message, url.Values{})
		if err != nil {
			logrus.Warningf("Twitter notification failure: %v", err)
		} else {
			logrus.Info("Twitter: Sent!")
			logrus.Info("https://twitter.com/%s/status/%s", tweet.User.Name, tweet.IdStr)
		}
	}

	return nil

}
