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
)

type Auth struct {
	hackmdToken string
	hackmdID    string

	//discordToken   string
	//discordChannel string
	//
	//twitterApiKey            string
	//twitterApiKeySecret      string
	//twitterBearerToken       string
	//twitterAccessToken       string
	//twitterAccessTokenSecret string
	//
	//mastodonClientID     string
	//mastodonClientSecret string
	//mastodonAccessToken  string
	//mastodonUsername     string
	//mastodonPassword     string
	//mastodonServer       string
}

var (
	auth     = &Auth{}
	registry = []*EnvVar{
		{
			Name:        "LIVE_HACKMD_TOKEN",
			Value:       "",
			Destination: &auth.hackmdToken,
			Required:    true,
		},
		{
			Name:        "LIVE_HACKMD_ID",
			Value:       "",
			Destination: &auth.hackmdID,
			Required:    true,
		},
	}
)

type EnvVar struct {
	Name        string
	Value       string
	Destination *string
	Required    bool
}

func LoadEnvironmentalVariables() error {
	for _, v := range registry {
		v.Value = os.Getenv(v.Name)
		if v.Required && v.Value == "" {
			return fmt.Errorf("empty env var [%s]", v.Name)
		}
		*v.Destination = v.Value
	}
	return nil
}
