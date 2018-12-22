package uptimerobot

import (
	"fmt"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	client = New("")

	os.Exit(m.Run())
}

func TestClient_GetMonitors(t *testing.T) {
	m, err := client.GetMonitors()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(m)
}
