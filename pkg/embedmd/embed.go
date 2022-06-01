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
	"os"
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

func ReadVFile(v interface{}, path string) error {
	var err error
	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("file read error: %s: %v", path, err)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %s: %v", path, err)
	}
	err = Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %v", err)
	}
	return nil
}

func RecordVFile(v interface{}, path string) error {
	var dataToWrite []byte
	var err error
	_, err = os.Stat(path)
	if err != nil {
		// File does not exist
		dataToWrite, err = RecordV([]byte(""), v)
		if err != nil {
			return fmt.Errorf("unable to RecordV: %v", err)
		}
	} else {
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read path: %s: %v", path, err)
		}
		dataToWrite, err = RecordV(fileBytes, v)
		if err != nil {
			return fmt.Errorf("unable to RecordV: %v", err)
		}
	}
	return ioutil.WriteFile(path, dataToWrite, DefaultPermission)
}

// RecordV will replace an embed source with a new one (if it exists), otherwise will append to the end
func RecordV(data []byte, v interface{}) ([]byte, error) {
	dataStr := string(data)
	rawBytes, err := Marshal(v)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to Marshal: %v", err)
	}
	if !strings.Contains(dataStr, EmbedDelim) {
		return embed(data, rawBytes)
	}
	return replace(data, rawBytes)
}

func embed(data, raw []byte) ([]byte, error) {
	dataStr := string(data)
	var str string
	str += "\n"
	str += "\n"
	str += EmbedStart
	str += string(raw)
	str += EmbedStop
	str += "\n"
	str += "\n"
	combine := dataStr + str
	return []byte(combine), nil
}

func replace(data, raw []byte) ([]byte, error) {
	dataStr := string(data)
	spl := strings.Split(dataStr, EmbedDelim)
	if len(spl) != 3 {
		return []byte(""), fmt.Errorf("invalid embed: not 3: %d", len(spl))
	}
	spl[1] = EmbedDelim + string(raw) + EmbedDelim
	str := strings.Join(spl, "")
	return []byte(str), nil
}

// MarshalMarkdown will cast any interface into a raw embed string (this will include the markdown syntax).
func MarshalMarkdown(v interface{}) ([]byte, error) {
	e, err := MarshalMarkdown(v)
	if err != nil {
		return []byte(""), err
	}
	var str string
	str += "\n"
	str += "\n"
	str += EmbedStart
	str += string(e)
	str += EmbedStop
	str += "\n"
	str += "\n"
	return []byte(str), nil
}

// Marshal will cast any interface into a raw embed string (raw, no syntax).
func Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte(""), fmt.Errorf("unable json marshal object: %v", err)
	}
	e := &E{
		Data:    data,
		Updated: time.Now(),
	}
	jbytes, err := json.Marshal(e)
	if err != nil {
		return []byte(""), fmt.Errorf("internal marshal failure: %v", err)
	}
	str := base64.StdEncoding.EncodeToString(jbytes)
	return []byte(str), nil
}

// Unmarshal will cast the embedded data onto an interface{]
func Unmarshal(data []byte, v interface{}) error {
	raw := &bytes.Buffer{}
	err := read(data, raw)
	if err != nil {
		return fmt.Errorf("unable to read: %v", err)
	}
	return json.Unmarshal(raw.Bytes(), v)
}

// read will read the raw (embedded data) directly from a data source (if it exists).
func read(data []byte, raw *bytes.Buffer) error {
	rawStr, err := findRaw(data)
	if err != nil {
		return fmt.Errorf("unable to find raw: %v", err)
	}
	if rawStr == "" {
		return nil
	}
	jbytes, err := base64.StdEncoding.DecodeString(rawStr)
	if err != nil {
		return fmt.Errorf("unable to decode: %v", err)
	}
	e := &E{}
	err = json.Unmarshal(jbytes, e)
	if err != nil {
		return fmt.Errorf("unable to unmarshal e: %v", err)
	}
	raw.Write(e.Data)
	return nil
}

// findRaw will find the raw string in a file if it exists, otherwise empty string
func findRaw(data []byte) (string, error) {
	fileContent := string(data)
	return between(fileContent, EmbedStart, EmbedStop), nil
}

func between(str string, start string, stop string) (result string) {
	spl1 := strings.Split(str, start)
	if len(spl1) != 2 {
		return ""
	}
	spl2 := strings.Split(spl1[1], stop)
	if len(spl2) != 2 {
		return ""
	}
	return spl2[0]
}
