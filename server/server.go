package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/henilmalaviya/golw/env"
	"github.com/henilmalaviya/golw/game"
	"github.com/henilmalaviya/golw/util"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return !env.Get().WebSocketOriginCheck
	},
	EnableCompression: true,
}

var registry = NewCommandRegistry()

func HandleConnection(conn *websocket.Conn, gm *game.Manager) {
	defer conn.Close()

	logger := util.GetLogger()
	clientAddr := conn.RemoteAddr().String()
	logger.Info("WebSocket connection established", "client", clientAddr)

	observer := NewObserver(conn, gm)
	defer func() {
		observer.Close()
		logger.Info("WebSocket connection closed", "client", clientAddr)
	}()

	for {
		var msg IncomingMessage

		err := conn.ReadJSON(&msg)

		if err != nil {
			// Only break on websocket close errors
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Debug("WebSocket connection terminated", "client", clientAddr, "reason", err.Error())
				break
			}
			// If it's a JSON error, continue
			logger.Warn("Invalid message received", "client", clientAddr, "error", err.Error())
			observer.SendOutgoingMessage(ErrorUnknownCommand)
			continue
		}

		go observer.HandleIncomingMessage(msg)
	}

}

func WebsocketHandler(gm *game.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := util.GetLogger()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Failed to upgrade WebSocket connection", "error", err.Error(), "client", r.RemoteAddr)
			http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
			return
		}

		HandleConnection(conn, gm)
	}
}
