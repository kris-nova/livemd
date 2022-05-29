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
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	token := os.Getenv(EnvironmentalVariableToken)
	if token == "" {
		t.Errorf("Unable to read [%s] environmental variable. Empty.", EnvironmentalVariableToken)
		t.FailNow()
	}
	client := New(token)
	user, err := client.Me()
	if err != nil {
		t.Errorf("Unable to get /me: %v", err)
	}
	data, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("Unable json print /me: %v", err)
	}
	t.Logf(string(data))
}

func TestClientNotes(t *testing.T) {
	token := os.Getenv(EnvironmentalVariableToken)
	if token == "" {
		t.Errorf("Unable to read [%s] environmental variable. Empty.", EnvironmentalVariableToken)
		t.FailNow()
	}
	client := New(token)
	v, err := client.GetNotes()
	if err != nil {
		t.Errorf("Unable to get: %v", err)
	}
	data, err := json.Marshal(&v)
	if err != nil {
		t.Errorf("Unable json print: %v", err)
	}
	t.Logf(string(data))
}

const (
	WellKnownNoteID string = "FoFqCCmrRYiHr4jNm_ckwg"
)

func TestClientNote(t *testing.T) {
	token := os.Getenv(EnvironmentalVariableToken)
	if token == "" {
		t.Errorf("Unable to read [%s] environmental variable. Empty.", EnvironmentalVariableToken)
		t.FailNow()
	}
	client := New(token)
	v, err := client.GetNote(WellKnownNoteID)
	if err != nil {
		t.Errorf("Unable to get: %v", err)
	}
	data, err := json.Marshal(&v)
	if err != nil {
		t.Errorf("Unable json print: %v", err)
	}
	t.Logf(string(data))
}

func TestClientCreateDelete(t *testing.T) {
	token := os.Getenv(EnvironmentalVariableToken)
	if token == "" {
		t.Errorf("Unable to read [%s] environmental variable. Empty.", EnvironmentalVariableToken)
		t.FailNow()
	}
	client := New(token)

	note := &Note{
		Title:   "TEST NOTE FROM UNIT TESTS",
		Content: "# Beeps Boops",
	}
	v, err := client.CreateNote(note)
	if err != nil {
		t.Errorf("Unable to create: %v", err)
	}
	t.Logf("Created note: %s", v.ID)
	err = client.DeleteNote(v.ID)
	if err != nil {
		t.Errorf("Unable to delete: %v", err)
	}
	t.Logf("Deleted note: %s", v.ID)

}
