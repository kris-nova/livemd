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

package livemd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kris-nova/livemd/pkg"

	"github.com/kris-nova/livemd/pkg/hackmd"

	"github.com/kris-nova/livemd/pkg/embedmd"
)

const (
	DefaultMode = 0644
)

type Link struct {
	URL      *url.URL
	Title    string
	Markdown string
}

type LiveMD struct {
	// path is the coupling to a local file on disk.
	path string

	// TwitchID is the twitch id for the VOD
	//
	// This MUST be set for every livemd as we couple this
	// software to Twitch.
	//
	//https://www.twitch.tv/videos/{{id}}
	//https://dashboard.twitch.tv/u/{{channel}}/content/video-producer/edit/{{id}}
	TwitchID      string
	TwitchChannel string

	Title       string
	Notify      string
	Description string

	// Subsystems
	Links   []*Link
	Twitter string
}

func New(path, twitchChannel, twitchID string) *LiveMD {
	return &LiveMD{
		path:          path,
		TwitchID:      twitchID,
		TwitchChannel: twitchChannel,
	}
}

func (l *LiveMD) TwitchVideoPage() string {
	return fmt.Sprintf("https://www.twitch.tv/videos/%s", l.TwitchID)
}

func (l *LiveMD) TwitchEditPage() string {
	return fmt.Sprintf("https://dashboard.twitch.tv/u/%s/content/video-producer/edit/%s", l.TwitchChannel, l.TwitchID)
}

func (l *LiveMD) AddLink(title, rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	l.Links = append(l.Links, &Link{
		URL:      u,
		Title:    title,
		Markdown: fmt.Sprintf("[%s](%s)", title, u.String()),
	})
	return nil
}

//// Write will first render the markdown, and write the result to the configured path
//func (l *LiveMD) Read() error {
//	rawData, err := ioutil.ReadFile(l.path)
//	if err != nil {
//		return fmt.Errorf("unable to read file: %s: %v", l.path, err)
//	}
//	return embedmd.Unmarshal(rawData, l)
//}

// Write will first render the markdown, and write the result to the configured path
func (l *LiveMD) Write() error {
	rawMarkdown, err := l.Render()
	if err != nil {
		return fmt.Errorf("unable to render markdown: %v", err)
	}
	return ioutil.WriteFile(l.path, rawMarkdown, DefaultMode)
}

// CoalesceDateName will turn a title and today's date into a unix friendly
// name that can be used for file creation, and archiving.
//
// h/t @mrxinu
func (l *LiveMD) CoalesceDateName() string {
	now := time.Now()
	formatted := now.Format("2006-01-02")
	title := strings.ToLower(l.Title)
	title = strings.ReplaceAll(title, " ", "_")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "!", "")
	return formatted + "__" + title + ".md"
}

// Read will read from disk exactly, with no mutations.
func (l *LiveMD) Read() ([]byte, error) {
	// Read the raw markdown from the filesystem
	readBytes, err := ioutil.ReadFile(l.path)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to read markdown: %v", err)
	}
	return readBytes, nil
}

// Render will render the markdown and return the content.
// Render will (by design) always template!
func (l *LiveMD) Render() ([]byte, error) {
	var rawMarkdown []byte
	// Build the raw markdown from the template
	tpl := template.New(l.path)
	tpl, err := tpl.Parse(pkg.MarkdownTemplate)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to parse template: %v", err)
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, l)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to execute template: %v", err)
	}
	rawMarkdown = buf.Bytes()
	// RecordV will record *l directly into rawMarkdown
	return embedmd.RecordV(rawMarkdown, l)
}

// Load will attempt to load a *LiveMD from a path
func Load(path string) (*LiveMD, error) {
	x := &LiveMD{
		path: path,
	}
	data, err := x.Read()
	if err != nil {
		return nil, err
	}
	err = embedmd.Unmarshal(data, x)
	return x, err
}

func ToNote(l *LiveMD, id string) (*hackmd.Note, error) {
	data, err := l.Render()
	if err != nil {
		return nil, err
	}
	note := &hackmd.Note{
		Title:   l.Title,
		ID:      id,
		Content: string(data),
	}
	return note, nil
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
	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("unable to move file. ensure directory exists for destination: %s: %v", dst, err)
	}
	return nil
}
