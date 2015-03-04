package content

import (
	"github.com/gorilla/websocket"
	"io"
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

type WebSocketConnInterface interface {
	NextReader() (int, io.Reader, error)
	Close() error
	GetKillChan() chan bool
}

type MockWebSocketConn struct {
	KillChan        chan bool
	MockGetKillChan func() chan bool
	MockNextReader  func() (int, io.Reader, error)
	MockClose       func() error
}

func (mwsc *MockWebSocketConn) NextReader() (int, io.Reader, error) {
	if mwsc.MockNextReader != nil {
		return mwsc.MockNextReader()
	}
	return 0, nil, nil
}
func (mwsc *MockWebSocketConn) Close() error {
	if mwsc.MockClose != nil {
		return mwsc.MockClose()
	}
	return nil
}
func (mwsc *MockWebSocketConn) GetKillChan() chan bool {
	if mwsc.MockGetKillChan != nil {
		return mwsc.MockGetKillChan()
	}
	return mwsc.KillChan
}

// NoOpReadLoop - If you don't care about reading from a websocket, drop it on the floor.
func NoOpReadLoop(wsConn WebSocketConnInterface, stopReading chan bool) {
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
				wsConn.GetKillChan() <- true
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

func (wsc *WebSocketConn) GetKillChan() chan bool {
	return wsc.KillChan
}

// NewWebSocketConn - create a new wrapper websocketconn
func NewWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	return &WebSocketConn{
		conn,
		make(chan bool),
	}
}
