package server

import (
	"encoding/json"

	"github.com/henilmalaviya/golw/util"
	"github.com/tidwall/gjson"
)

type CommandHandler func(data gjson.Result, observer *Observer, wc chan<- OutgoingMessage)

type CommandRegistry struct {
	handlers map[Command]CommandHandler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[Command]CommandHandler),
	}
}

func (r *CommandRegistry) Register(command Command, handler CommandHandler) {
	r.handlers[command] = handler
}

func (r *CommandRegistry) Handle(command Command, data MessageData, observer *Observer) {
	logger := util.GetLogger()
	handler, exists := r.handlers[command]
	if !exists {
		logger.Warn("Unknown command received", "command", string(command))
		observer.SendOutgoingMessage(ErrorUnknownCommand)
		return
	}

	marshalData, err := json.Marshal(data)

	if err != nil {
		logger.Error("Failed to marshal command data", "command", string(command), "error", err.Error())
		observer.SendOutgoingMessage(ErrorUnknownCommand)
		return
	}

	parsedData := gjson.ParseBytes(marshalData)

	logger.Debug("Executing command", "command", string(command))

	ch := NewOutgoingMessageChannel()
	go handler(parsedData, observer, ch)
	for msg := range ch {
		if err := observer.SendOutgoingMessage(msg); err != nil {
			logger.Error("Failed to send message to client", "error", err.Error())
			return
		}
	}
	close(ch)
}
