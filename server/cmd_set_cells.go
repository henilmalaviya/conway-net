package server

import (
	"github.com/henilmalaviya/golw/util"
	"github.com/tidwall/gjson"
)

func CommandSetCellsHandler(data gjson.Result, observer *Observer, wc chan<- OutgoingMessage) {
	logger := util.GetLogger()

	cellsArray, ok := util.GetCellsArrayFromData(data, "cells")
	if !ok {
		logger.Warn("Invalid cells data received in set_cells command")
		wc <- NewOutgoingErrorMessage("invalid cells data")
		return
	}

	gm := observer.Manager.GetGame()
	logger.Info("Setting cells", "count", len(cellsArray))

	for _, cell := range cellsArray {
		gm.GetGrid().SetCell(cell.X, cell.Y)
	}

	wc <- OkMessage
}

func init() {
	registry.Register(CommandSetCells, CommandSetCellsHandler)
}
