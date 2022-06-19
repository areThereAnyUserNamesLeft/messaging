package message

import (
	"encoding/json"
	"fmt"
	"log"
)

type Header struct {
	UUID       string
	Role       string
	Recipients []string
}

type Body struct {
	Content string
}

type Message struct {
	Header Header
	Body   Body
}

// NewMessage creates a message as a collection of structs
func NewMessage(role string, messContents string, ID string, recipients []string) Message {
	return Message{
		Header: Header{UUID: ID, Role: role, Recipients: recipients},
		Body:   Body{Content: messContents}}
}

//Makes message into JSON string
func makeJSON(message Message) string {
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("JSON Marshall Error: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

// I decided to seperate these here to make it clear when calling them

//Message Flavoring - List
func ListMessage(content string, ID string) string {
	return fmt.Sprintf("%s", makeJSON(NewMessage("List", content, ID, nil)))
}

//Message Flavoring - Identity
func IdentityMessage(content string, ID string) string {
	return fmt.Sprintf("%s", makeJSON(NewMessage("Identity", content, ID, nil)))
}

//Message Flavoring - Relay
func RelayMessage(content string, ID string, recipients []string) string {
	return fmt.Sprintf("%s", makeJSON(NewMessage("Relay", content, ID, recipients)))
}
