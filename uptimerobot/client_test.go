package uptimerobot

import (
	"os"
	"testing"

	"github.com/kr/pretty"
)

var client *Client

func TestMain(m *testing.M) {
	client = New("u505705-4387f0029e8cc3caa22b0fb7")

	os.Exit(m.Run())
}

var id int

func TestClient_NewMonitor(t *testing.T) {
	var err error
	id, err = client.NewMonitor(Monitor{
		URL:          "https://google.com",
		FriendlyName: "Google Test",
		Type:         MonitorTypeHTTP,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetMonitors(t *testing.T) {
	monitors, err := client.GetMonitors()
	if err != nil {
		t.Error(err)
	}

	var found bool
	for _, m := range monitors {
		if m.ID == id {
			pretty.Println(m)
			found = true
		}
	}
	if !found {
		t.Error("newly created monitor", id, "not present")
	}
}

func TestClient_DeleteMonitor(t *testing.T) {
	err := client.DeleteMonitor(id)
	if err != nil {
		t.Error(err)
	}
}
