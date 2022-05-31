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

package embedmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

const (
	EmbedDelim string = ":embed:"
	EmbedStart string = "```\n:embed:"
	EmbedStop  string = ":embed:\n```"
)

const (
	DefaultPermission = 0655
)

type E struct {
	Updated time.Time
	Data    []byte
}

func Write(data []byte, path string) error {
	e := &E{
		Data:    data,
		Updated: time.Now(),
	}
	jbytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("internal marshal failure: %v", err)
	}
	str := base64.StdEncoding.EncodeToString(jbytes)
	raw, err := findRaw(path)
	if err != nil {
		return fmt.Errorf("internal marshal failure: %v", err)
	}
	if raw == "" {
		return embed(str, path)
	} else {
		return replace(raw, path)
	}
	return nil
}

func WriteV(v interface{}, path string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("unable json marshal object: %v", err)
	}
	return Write(data, path)
}

func ReadV(v interface{}, path string) error {
	buf := &bytes.Buffer{}
	err := Read(buf, path)
	if err != nil {
		return fmt.Errorf("unable to read: %v", err)
	}
	return json.Unmarshal(buf.Bytes(), v)
}

func Read(buf *bytes.Buffer, path string) error {
	raw, err := findRaw(path)
	if err != nil {
		return fmt.Errorf("unable to find raw: %v", err)
	}
	if raw == "" {
		return nil
	}
	jbytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("unable to decode: %v", err)
	}
	e := &E{}
	err = json.Unmarshal(jbytes, e)
	if err != nil {
		return fmt.Errorf("unable to unmarshal e: %v", err)
	}
	buf.Write(e.Data)
	return nil
}

// findRaw will find the raw string in a file if it exists, otherwise empty string
func findRaw(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	fileContent := string(data)
	return between(fileContent, EmbedStart, EmbedStop), nil
}

func between(str string, start string, stop string) (result string) {
	return strings.TrimLeft(strings.TrimRight(str, stop), start)
}

func replace(raw string, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fileContent := string(data)
	if !strings.Contains(fileContent, EmbedDelim) {
		return fmt.Errorf("missing embed in file")
	}
	spl := strings.Split(fileContent, EmbedDelim)
	if len(spl) != 3 {
		// 1. before
		// 2. raw
		// 3. after
		return fmt.Errorf("invalid embed: not 3: %d", len(spl))
	}
	spl[1] = raw
	str := strings.Join(spl, "")
	return ioutil.WriteFile(path, []byte(str), DefaultPermission)
}

// embed will append the raw content to a file
func embed(raw string, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fileContent := string(data)
	// ----------------
	str := fileContent
	str += "\n"
	str += "\n"
	str += EmbedStart
	str += raw
	str += EmbedStop
	str += "\n"
	str += "\n"
	// ----------------
	return ioutil.WriteFile(path, []byte(str), DefaultPermission)
}
