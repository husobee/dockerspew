package content

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// NewWebSocketUpgrader
func NewWebSocketUpgrader(readBufferSize, writeBufferSize int) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

// NoOpReadLoop - If you don't care about reading from a websocket, drop it on the floor.
func NoOpReadLoop(wsConn *WebSocketConn, stopReading chan bool) {
OUTER:
	for {
		select {
		case <-stopReading:
			close(stopReading)
			return
		default:
			if _, _, err := wsConn.NextReader(); err != nil {
				log.Println("[ERROR] Failed to grab next reader, err=", err.Error())
				wsConn.Close()
				wsConn.KillChan <- true
				break OUTER
			}
		}
	}
}

// WebSocketConn - wrap gorilla websocket with a
type WebSocketConn struct {
	*websocket.Conn
	KillChan chan bool
}

// NewWebSocketConn - create a new wrapper websocketconn
func NewWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	return &WebSocketConn{
		conn,
		make(chan bool),
	}
}
