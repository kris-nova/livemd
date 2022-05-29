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
