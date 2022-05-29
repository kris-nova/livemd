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

	"github.com/kris-nova/live/pkg/hackmd"
)

const (
	LiveMDPerm = 0655
)

type LiveMD struct {
	Title     string
	I         int
	YouTubeID string
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

func FromRaw(data []byte) (*LiveMD, error) {
	x := &LiveMD{
		Data: data,
	}

	return x, nil
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
	return buf.Bytes(), nil
}
