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
	"testing"
)

type TestObject struct {
	Name    string
	Content string
	Number  int
}

func TestValidExisting(t *testing.T) {
	v := &TestObject{
		Name:   "barnaby",
		Number: 7,
	}
	path := "testdata/valid_existing.md"
	err := WriteV(v, path)
	if err != nil {
		t.Errorf("unable to WriteV(): %v", err)
	}
	x := &TestObject{}
	err = ReadV(x, "testdata/valid_existing.md")
	if err != nil {
		t.Errorf("unable to Read(): %v", err)
	}
	if x.Name != "barnaby" {
		t.Errorf("invalid read: expecing matching name")
	}
}
