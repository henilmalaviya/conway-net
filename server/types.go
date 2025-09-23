package server

import (
	"encoding/json"
	"fmt"
)

type Command string

type Code string

type MessageData map[string]interface{}

const (
	CommandSetCells   Command = "set_cells"
	CommandClearCells Command = "clear_cells"
	CommandSync       Command = "sync"
	CommandObserve    Command = "observe"
	CommandUnobserve  Command = "unobserve"
)

const (
	CodeOk           Code = "ok"
	CodeObserveOk    Code = "observe_ok"
	CodeSyncOk       Code = "sync_ok"
	CodeError        Code = "error"
	CodeObserveEvent Code = "observe_event"
)

type IncomingMessage struct {
	Command Command     `json:"command"`
	Data    MessageData `json:"data,omitempty"`
}

type OutgoingMessage struct {
	Code Code        `json:"code"`
	Data MessageData `json:"data,omitempty"`
}

func (o OutgoingMessage) String() string {
	if o.Data == nil {
		o.Data = make(MessageData)
	}
	jsonData, err := json.Marshal(o.Data)
	if err != nil || len(jsonData) == 0 {
		return fmt.Sprintf("%s;{}\r\n", o.Code)
	}
	return fmt.Sprintf("%s;%v\r\n", o.Code, string(jsonData))
}

func NewOutgoingMessage(code Code, data MessageData) OutgoingMessage {
	return OutgoingMessage{
		Code: code,
		Data: data,
	}
}

func NewOutgoingErrorMessage(err string) OutgoingMessage {
	return NewOutgoingMessage(CodeError, MessageData{"error": err})
}

func NewOutgoingMessageChannel() chan OutgoingMessage {
	return make(chan OutgoingMessage, 100)
}

var ErrorUnknownCommand = NewOutgoingErrorMessage("unknown command")
var ErrorInvalidData = NewOutgoingErrorMessage("invalid data")
var ErrorGameNotFound = NewOutgoingErrorMessage("game not found")

var OkMessage = NewOutgoingMessage(CodeOk, MessageData{})
