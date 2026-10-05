package connection

import (
	"testing"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func TestNewConnection(t *testing.T) {
	conn := NewConnection()

	if conn.Client != nil {
		t.Errorf("Client = %v, want nil", conn.Client)
	}
	if conn.UpdatesChannel == nil {
		t.Errorf("UpdatesChannel = nil, want non-nil")
	}
}

func TestSetMe(t *testing.T) {
	conn := NewConnection()

	me := conn.SetMe(&tdlib.User{Id: 42, FirstName: "Ada", LastName: "Lovelace"})

	if me.Id != 42 || me.FirstName != "Ada" || me.LastName != "Lovelace" {
		t.Errorf("SetMe() returned %+v, want {42 Ada Lovelace}", me)
	}

	got := conn.GetMe()
	if got != me {
		t.Errorf("GetMe() = %+v, want %+v", got, me)
	}
}
