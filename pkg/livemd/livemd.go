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
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/kris-nova/live/pkg"

	"github.com/kris-nova/live/pkg/hackmd"
)

const (
	LiveMDPerm = 0655
)

type LiveMD struct {
	Title     string
	I         int
	YouTubeID string
	TwitchID  string
	HackMDID  string
	Data      []byte
}

func New(title string) *LiveMD {
	return &LiveMD{
		Title: title,
		I:     I(),
	}
}

func FromFile(path string) (*LiveMD, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromRaw(data)
}

func FromHackMD(client *hackmd.Client, id string) (*LiveMD, error) {
	note, err := client.GetNote(id)
	if err != nil {
		return nil, err
	}
	return FromRaw([]byte(note.Content))
	return nil, nil
}

const (
	DataStartDelim string = "data:\n" + "```json\n"
	DataStopDelim  string = "\n```\n" + "data:\n"
)

func FromRaw(data []byte) (*LiveMD, error) {
	x := &LiveMD{}
	rawBytes, err := findRaw(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawBytes, x)
	if err != nil {
		logrus.Warnf(string(rawBytes))
		return nil, fmt.Errorf("unable to unmarshal raw: %v", err)
	}
	x.Data = data // Always overwrite the data at the end
	return x, nil
}

// findRaw will find the embedded raw data in the content
func findRaw(data []byte) ([]byte, error) {
	str := string(data)
	spl := strings.Split(str, DataStartDelim)
	if len(spl) != 2 {
		return nil, fmt.Errorf("invalid DataStartDelim")
	}
	spll := strings.Split(spl[1], DataStopDelim)
	if len(spll) != 2 {
		return nil, fmt.Errorf("invalid DataStopDelim")
	}
	raw := spll[0]
	return []byte(raw), nil
}

func (x *LiveMD) Write(path string) error {
	md, err := x.Markdown()
	if err != nil {
		return fmt.Errorf("unable to format markdown: %v", err)
	}
	return ioutil.WriteFile(path, md, LiveMDPerm)
}

// Markdown is a deterministic function based on the runtime configuration of *LiveMD
// Markdown will template the live.md file in /pkg
//
// Markdown will NOT write from disk.
func (x *LiveMD) Markdown() ([]byte, error) {
	tpl := template.New(x.Title)
	tpl, err := tpl.Parse(pkg.MarkdownTemplate)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to parse template: %v", err)
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, x)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to execute template: %v", err)
	}
	rawBytes, err := findRaw(buf.Bytes())
	if err != nil {
		return []byte(""), fmt.Errorf("unable to find raw data: %v", err)
	}
	x.Data = []byte("") // Unset data
	newBytes, err := json.MarshalIndent(x, "", "   ")
	if err != nil {
		return []byte(""), fmt.Errorf("unable to marshal new bytes: %v", err)
	}
	str := strings.Replace(buf.String(), string(rawBytes), string(newBytes), 1)
	return []byte(str), nil
}

// HackMDNote will convert a *LiveMD to a *hackmd.Note with an optional ID (can be empty)
func (x *LiveMD) HackMDNote(id string) (*hackmd.Note, error) {
	note := &hackmd.Note{
		ID:      id,
		Title:   x.Title,
		Content: string(x.Data),
	}
	return note, nil
}
