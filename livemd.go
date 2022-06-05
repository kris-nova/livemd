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

	"github.com/kris-nova/livemd/pkg"

	"github.com/kris-nova/livemd/pkg/embedmd"
)

const (
	DefaultMode = 0644
)

type LiveMD struct {
	// path is the coupling to a local file on disk.
	path string

	// Title is the title of the show
	Title       string
	Notify      string
	Description string
}

func New(path string) *LiveMD {
	return &LiveMD{
		path: path,
	}
}

// Write will first render the markdown, and write the result to the configured path
func (l *LiveMD) Read() error {
	rawData, err := ioutil.ReadFile(l.path)
	if err != nil {
		return fmt.Errorf("unable to read file: %s: %v", l.path, err)
	}
	return embedmd.Unmarshal(rawData, l)
}

// Write will first render the markdown, and write the result to the configured path
func (l *LiveMD) Write() error {
	rawMarkdown, err := l.Render()
	if err != nil {
		return fmt.Errorf("unable to render markdown: %v", err)
	}
	return ioutil.WriteFile(l.path, rawMarkdown, DefaultMode)
}

// Render will render the markdown and return the content
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
	x := New(path)
	err := x.Read()
	return x, err
}
