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
	"io/ioutil"
	"testing"
)

type TestObject struct {
	Name    string
	Content string
	Number  int
}

func TestFindRaw(t *testing.T) {
	needle := "TestDataRaw"
	data, err := ioutil.ReadFile("testdata/simple_raw.md")
	if err != nil {
		t.Errorf("unable to read test: %v", err)
	}
	find, err := findRaw(data)
	if err != nil {
		t.Errorf("unable to findRaw: %v", err)
	}
	if find != needle {
		t.Errorf("expecting: %s, found: %s", needle, find)
	}
}

func TestValidExisting(t *testing.T) {
	v := &TestObject{
		Name:   "barnaby",
		Number: 7,
	}
	path := "testdata/valid_existing.md"
	err := RecordVFile(v, path)
	if err != nil {
		t.Errorf("unable to WriteV(): %v", err)
	}
	x := &TestObject{}
	err = ReadVFile(x, "testdata/valid_existing.md")
	if err != nil {
		t.Errorf("unable to Read(): %v", err)
	}
	if x.Name != "barnaby" {
		t.Errorf("invalid read: expecing matching name")
	}
	if x.Number != 7 {
		t.Errorf("invalid read: expecing matching number")

	}
}
