package server

import (
	"github.com/henilmalaviya/golw/util"
	"github.com/tidwall/gjson"
)

func CommandSyncHandler(data gjson.Result, observer *Observer, wc chan<- OutgoingMessage) {
	logger := util.GetLogger()

	gm := observer.Manager.GetGame()
	gr := gm.GetGrid()

	bounds, ok := util.GetBoundsFromData(data, "bounds")

	if !ok {
		logger.Warn("Sync command received without valid bounds data")
		wc <- NewOutgoingErrorMessage("bounds not provided")
		return
	}

	liveCells := gr.GetLiveCellCoordinates()
	logger.Debug("Syncing grid state", "bounds", bounds.ToNestedArray(), "live_cells_count", len(liveCells))

	wc <- NewOutgoingMessage(CodeSyncOk, MessageData{"cells": liveCells, "bounds": bounds.ToNestedArray(), "stats": observer.Manager.GetStats()})
}

func init() {
	registry.Register(CommandSync, CommandSyncHandler)
}
