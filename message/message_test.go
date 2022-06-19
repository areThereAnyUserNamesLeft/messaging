package message_test

import (
	"messaging/message"
	"strings"
	"testing"
)

func TestMessage(t *testing.T) {
	r := message.RelayMessage("This is my content!", "id-relay", nil)
	if strings.Index(r, "Relay") == -1 {
		t.Error("Role not in message", r)
	}
	if strings.Index(r, "This is my content!") == -1 {
		t.Error("Content not in message", r)
	}
	l := message.ListMessage("This is my content!", "id-list")
	if strings.Index(l, "List") == -1 {
		t.Error("Role not in message", l)
	}
	if strings.Index(l, "This is my content!") == -1 {
		t.Error("Content not in message", l)
	}
	i := message.IdentityMessage("This is my content!", "id-identity")
	if strings.Index(i, "Identity") == -1 {
		t.Error("Role not in message", i)
	}
	if strings.Index(i, "This is my content!") == -1 {
		t.Error("Content not in message", i)
	}
}
