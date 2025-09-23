package server

import (
	"github.com/henilmalaviya/golw/util"
	"github.com/tidwall/gjson"
)

func CommandUnobserveHandler(data gjson.Result, observer *Observer, wc chan<- OutgoingMessage) {
	logger := util.GetLogger()

	if observer.gridObserver == nil {
		logger.Warn("Unobserve command received but client is not observing any grid")
		wc <- NewOutgoingErrorMessage("not observing any grid")
		return
	}

	gm := observer.Manager.GetGame()
	gr := gm.GetGrid()

	gr.RemoveObserver(observer.gridObserver)
	observer.gridObserver = nil

	logger.Info("Client stopped observing grid")
	wc <- OkMessage
}

func init() {
	registry.Register(CommandUnobserve, CommandUnobserveHandler)
}
