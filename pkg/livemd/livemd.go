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

	"github.com/kris-nova/live/pkg"

	"github.com/kris-nova/live/pkg/embedmd"

	"github.com/kris-nova/live/pkg/hackmd"
)

type LiveMD struct {
	content     []byte
	Initialized bool
	Title       string
	I           int
	YouTubeID   string
	TwitchID    string
	HackMDID    string
}

func New(title string) *LiveMD {
	return &LiveMD{
		Title: title,
		I:     I(),
	}
}

// FromLocal will return a new *LiveMD from a local path
func FromLocal(path string) (*LiveMD, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromRaw(data)
}

// FromRemote will fetch a *LiveMD from a remote HackMD client
func FromRemote(client *hackmd.Client, id string) (*LiveMD, error) {
	note, err := client.GetNote(id)
	if err != nil {
		return nil, err
	}
	return FromRaw([]byte(note.Content))
	return nil, nil
}

// FromRaw is a determinstic function that will return a *LiveMD from a raw data source
func FromRaw(raw []byte) (*LiveMD, error) {
	x := &LiveMD{}
	err := embedmd.Unmarshal(raw, x)

	return x, err
}

func (x *LiveMD) Write(path string) error {
	return ioutil.WriteFile(path, x.content, embedmd.DefaultPermission)
}

// Data will return the raw file data calculated for a *LiveMD
func (x *LiveMD) Data() ([]byte, error) {
	var preData []byte
	if len(x.content) > 0 {
		// Preexisting
		preData = x.content
	} else {
		// New
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
		preData = buf.Bytes()
	}
	x.content = preData // Always update the content before writing
	return embedmd.RecordV(preData, x)
}

func (x *LiveMD) SyncRaw(data []byte) {
	x.content = data
}

// Sync will update the embedded data from the actual content on disk
func (x *LiveMD) Sync(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	ret, err := embedmd.RecordV(data, x)
	if err != nil {
		return err
	}
	x.content = ret
	return nil
}

// ToHackMD will convert a *LiveMD to a *hackmd.Note with an optional ID (can be empty)
func (x *LiveMD) ToHackMD(id string) *hackmd.Note {
	return &hackmd.Note{
		ID:      id,
		Title:   x.Title,
		Content: string(x.content),
	}
}
