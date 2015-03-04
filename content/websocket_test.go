package content

import (
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestNewWebSocketUpgrader - create a new websocket upgrader
func TestNewWebSocketUpgrader(t *testing.T) {
	Convey("create a new websocket upgrader", t, func() {
		u := NewWebSocketUpgrader(1024, 2048)
		So(u.ReadBufferSize, ShouldEqual, 1024)
		So(u.WriteBufferSize, ShouldEqual, 2048)
	})
}

// TestNoOpReadLoop - Drop WebSocket Reads
func TestNoOpReadLoop(t *testing.T) {
	Convey("create a new websocket connection", t, func() {
		NewWebSocketConn(&websocket.Conn{})
		c := &MockWebSocketConn{}
		s := make(chan bool)
		ms := make(chan bool)
		mc := &MockWebSocketConn{
			KillChan: ms,
		}
		Convey("start the NoOpReadLoop", func() {
			go NoOpReadLoop(c, s)
		})
		Convey("start the NoOpReadLoop with mock", func() {
			go NoOpReadLoop(mc, ms)
			ms <- true
		})
	})
}
