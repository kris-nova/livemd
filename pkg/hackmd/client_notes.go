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

package hackmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LastChangeUser struct {
}

type Note struct {
	ID              string `json:"id,omitempty"`
	Title           string `json:"title,omitempty"`
	PublishType     string `json:"publishType,omitempty"`
	Permalink       string `json:"permalink,omitempty"`
	ShortID         string `json:"shortId,omitempty"`
	LastChangeUser  User
	UserPath        string `json:"userPath,omitempty"`
	TeamPath        string `json:"teamPath,omitempty"`
	ReadPermission  string `json:"readPermission,omitempty"`
	WritePermission string `json:"writePermission,omitempty"`
}

func (c *Client) Notes() ([]*Note, error) {
	resp, err := c.GET("notes")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Response %d %s\n", resp.StatusCode, string(data))
	}
	var v []*Note
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("unable to JSON unmarshal body: %v", err)
	}
	return v, nil
}

func (c *Client) Note(id string) (*Note, error) {
	resp, err := c.GET(fmt.Sprintf("notes/%s", id))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Response %d %s\n", resp.StatusCode, string(data))
	}
	var v *Note
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("unable to JSON unmarshal body: %v", err)
	}
	return v, nil
}
