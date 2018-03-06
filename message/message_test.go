package message_test

import (
	"fmt"
	"strings"
	"testing"
	"unity3dAssignment/message"
)

func TestMessage(t *testing.T) {
	r := message.OutRelayMess("This is my content!")
	if strings.Index(r, "Relay") == -1 {
		t.Error("Role not in message", r)
	}
	if strings.Index(r, "This is my content!") == -1 {
		t.Error("Content not in message", r)
	}
	l := message.OutListMess("This is my content!")
	if strings.Index(l, "List") == -1 {
		t.Error("Role not in message", l)
	}
	if strings.Index(l, "This is my content!") == -1 {
		t.Error("Content not in message", l)
	}
	i := message.OutIdentityMess("This is my content!")
	if strings.Index(i, "Identity") == -1 {
		t.Error("Role not in message", i)
	}
	if strings.Index(i, "This is my content!") == -1 {
		t.Error("Content not in message", i)
	}
	fmt.Println(r)
	fmt.Println(l)
	fmt.Println(i)
}
