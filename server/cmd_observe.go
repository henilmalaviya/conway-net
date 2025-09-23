package server

import (
	"fmt"

	"github.com/henilmalaviya/gol/grid"
	"github.com/henilmalaviya/golw/env"
	"github.com/henilmalaviya/golw/util"
	"github.com/tidwall/gjson"
)

func cellSliceToIntSlice(cells []grid.Cell) [][]int {
	parsedCells := make([][]int, len(cells))
	for i, cell := range cells {
		parsedCells[i] = []int{cell.X, cell.Y}
	}
	return parsedCells
}

func CommandObserveHandler(data gjson.Result, observer *Observer, wc chan<- OutgoingMessage) {
	logger := util.GetLogger()

	gm := observer.Manager.GetGame()
	gr := gm.GetGrid()

	var obs *grid.RegionObserver

	bounds, ok := util.GetBoundsFromData(data, "bounds")
	if !ok {
		logger.Warn("Invalid bounds data received in observe command")
		wc <- NewOutgoingErrorMessage("invalid bounds data")
		return
	}

	if bounds.Width() <= 0 || bounds.Height() <= 0 {
		logger.Warn("Invalid bounds dimensions", "width", bounds.Width(), "height", bounds.Height())
		wc <- NewOutgoingErrorMessage("bounds must have positive width and height")
		return
	}

	boundDiagonalLength := util.DiagonalLength(bounds)
	if boundDiagonalLength > float64(env.Get().MaxObserveRegionSize) {
		logger.Warn("Bounds exceed maximum region size", "diagonal", boundDiagonalLength, "max", env.Get().MaxObserveRegionSize)
		wc <- NewOutgoingErrorMessage(fmt.Sprintf("bounds diagonal length exceeds maximum allowed (%d)", env.Get().MaxObserveRegionSize))
		return
	}

	logger.Debug("Setting up region observer", "width", bounds.Width(), "height", bounds.Height())

	var updateFunc func(event grid.ObserverEvent) = func(event grid.ObserverEvent) {
		switch e := event.(type) {
		case grid.SetCellsObserverEvent:
			cells := e.Data()

			parsedCells := cellSliceToIntSlice(cells)

			wc <- NewOutgoingMessage(CodeObserveEvent, MessageData{
				"event": e.Type(),
				"data": map[string][][]int{
					"cells": parsedCells,
				},
			})
		case grid.SetCellObserverEvent:
			cell := e.Data()

			wc <- NewOutgoingMessage(CodeObserveEvent, MessageData{
				"event": e.Type(),
				"data": map[string][2]int{
					"cell": {cell.X, cell.Y},
				},
			})
		case grid.ClearCellObserverEvent:
			cell := e.Data()

			wc <- NewOutgoingMessage(CodeObserveEvent, MessageData{
				"event": e.Type(),
				"data": map[string][2]int{
					"cell": {cell.X, cell.Y},
				},
			})
		case grid.TickObserverEvent:
			bornCells, diedCells := e.Data()

			if len(bornCells) == 0 && len(diedCells) == 0 {
				return
			}

			parsedBornCells := cellSliceToIntSlice(bornCells)
			parsedDiedCells := cellSliceToIntSlice(diedCells)

			wc <- NewOutgoingMessage(CodeObserveEvent, MessageData{
				"event": e.Type(),
				"data": map[string][][]int{
					"bornCells": parsedBornCells,
					"diedCells": parsedDiedCells,
				},
			})
		}
	}

	if observer.gridObserver != nil {
		observer.gridObserver.SetRegion(bounds)
	} else {
		obs = grid.NewRegionObserver(bounds, updateFunc)
		observer.gridObserver = obs
		gr.AddObserver(observer.gridObserver)
	}

	wc <- NewOutgoingMessage(CodeObserveOk, MessageData{
		"bounds": bounds.ToNestedArray(),
	})
}

func init() {
	registry.Register(CommandObserve, CommandObserveHandler)
}
