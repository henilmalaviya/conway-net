package server

import (
	"compress/flate"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/henilmalaviya/gol/grid"
	"github.com/henilmalaviya/golw/game"
	"github.com/henilmalaviya/golw/util"
)

type Observer struct {
	Conn    *websocket.Conn
	Manager *game.Manager

	gridObserver *grid.RegionObserver

	connWriteMutex sync.Mutex
}

func (o *Observer) Close() {
	o.Conn.Close()
	gr := o.Manager.GetGame().GetGrid()

	if o.gridObserver != nil {
		gr.RemoveObserver(o.gridObserver)
	}

	o.gridObserver = nil
}

func NewObserver(conn *websocket.Conn, gm *game.Manager) *Observer {
	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(flate.BestSpeed)
	return &Observer{
		Conn:    conn,
		Manager: gm,
	}
}

func (o *Observer) HandleIncomingMessage(msg IncomingMessage) {
	logger := util.GetLogger()
	logger.Debug("Processing incoming message", "command", string(msg.Command))
	registry.Handle(msg.Command, msg.Data, o)
}

func (o *Observer) SendOutgoingMessage(msg OutgoingMessage) error {
	o.connWriteMutex.Lock()
	defer o.connWriteMutex.Unlock()
	return o.Conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}
