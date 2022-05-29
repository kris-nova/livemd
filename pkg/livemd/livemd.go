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
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/kris-nova/live/pkg/hackmd"
)

const (
	LiveMDPerm = 0655
)

type LiveMD struct {
	Title string
}

func New(title string) *LiveMD {
	return &LiveMD{
		Title: title,
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
	var livemd *LiveMD
	return livemd, nil
	// TODO Build from raw
}

func (x *LiveMD) Write(path string) error {
	md, err := x.Markdown()
	if err != nil {
		return fmt.Errorf("unable to format markdown: %v", err)
	}
	return ioutil.WriteFile(path, md, LiveMDPerm)
}

func (x *LiveMD) Markdown() ([]byte, error) {
	logrus.Warnf("MARKDOWN NOT SUPPORTED")
	return []byte(""), nil
}
