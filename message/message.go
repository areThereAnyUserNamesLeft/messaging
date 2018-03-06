package message

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
)

type Header struct {
	Uid  uint64
	Role string
}

type Body struct {
	Content string
}

type Message struct {
	Header Header
	Body   Body
}

// Creates a message as a collection of structs
func MkMess(role string, messContents string, Uid uint64) Message {
	mess := Message{
		Header: Header{Uid: Uid, Role: role},
		Body:   Body{Content: messContents}}
	return mess
}

// Lazy - Return here
func Uint64() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

//Makes message into JSON string
func mkJson(mess Message) string {
	data, err := json.Marshal(mess)
	if err != nil {
		log.Fatalf("JSON Marshall Error: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

// I decided to seperate these here to make it clear when calling them

//Message Flavoring - List
func ListMess(content string, Uid uint64) string {
	out := fmt.Sprintf("%s", mkJson(MkMess("List", content, Uid)))
	return out
}

//Message Flavoring - Identity
func IdentityMess(content string, Uid uint64) string {
	out := fmt.Sprintf("%s", mkJson(MkMess("Identity", content, Uid)))
	return out
}

//Message Flavoring - Relay
func RelayMess(content string, Uid uint64) string {
	out := fmt.Sprintf("%s", mkJson(MkMess("Relay", content, Uid)))
	return out
}
