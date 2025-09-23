package main

import (
	"net/http"
	"time"

	"github.com/henilmalaviya/golw/env"
	"github.com/henilmalaviya/golw/game"
	"github.com/henilmalaviya/golw/server"
	"github.com/henilmalaviya/golw/util"
)

func main() {
	logger := util.GetLogger()
	logger.Info("Starting Game of Life WebSocket server", "log_level", env.Get().LogLevel)

	gm := game.NewManager()
	logger.Info("Game manager initialized")

	saveManager := game.NewSaveManager(gm)

	if err := saveManager.LoadLatest(); err != nil {
		logger.Error("Failed to load latest snapshot", "error", err)
		// load a blinker pattern to start with
		gm.GetGame().GetGrid().SetCell(0, 1)
		gm.GetGame().GetGrid().SetCell(0, 0)
		gm.GetGame().GetGrid().SetCell(0, -1)

		logger.Info("Initialized game with default blinker pattern")

	} else {
		logger.Info("Loaded latest snapshot successfully")
		logger.Info("Current game stats", "stats", gm.GetStats())
	}

	saveManager.StartSaving()

	gm.Start(time.Millisecond * time.Duration(env.Get().TickSpeed))
	logger.Info("Game tick started", "interval", env.Get().TickSpeed)

	http.HandleFunc(env.Get().WSEndpoint, server.WebsocketHandler(gm))
	logger.Info("WebSocket endpoint registered", "endpoint", env.Get().WSEndpoint)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	logger.Info("Starting HTTP server", "port", env.Get().Port)
	if err := http.ListenAndServe(":"+env.Get().Port, nil); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
