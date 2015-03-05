package content

import (
	"errors"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"testing"
)

// TestNewWebSocketUpgrader - create a new websocket upgrader
func TestNewWebSocketUpgrader(t *testing.T) {
	Convey("create a new websocket upgrader", t, func() {
		u := NewWebSocketUpgrader(1024, 2048)
		So(u.ReadBufferSize, ShouldEqual, 1024)
		So(u.WriteBufferSize, ShouldEqual, 2048)
		So(u.CheckOrigin(nil), ShouldBeTrue)
	})
}

// TestNoOpReadLoop - Drop WebSocket Reads
func TestNoOpReadLoop(t *testing.T) {
	Convey("create a new websocket connection", t, func() {
		nwsc := NewWebSocketConn(&websocket.Conn{})
		kc := make(chan bool)
		nwsc.KillChan = kc
		So(nwsc.GetKillChan(), ShouldEqual, kc)

		close(kc)

		s := make(chan bool)
		ms := make(chan bool)
		ms1 := make(chan bool)
		readyFail := make(chan bool)

		c := &MockWebSocketConn{}
		mc := &MockWebSocketConn{
			KillChan: ms,
			MockNextReader: func() (int, io.Reader, error) {
				<-readyFail
				return 0, nil, errors.New("fail")
			},
		}
		mc1 := &MockWebSocketConn{}
		Convey("start the NoOpReadLoop", func() {
			go func() { s <- true }()
			go NoOpReadLoop(c, s)
		})
		Convey("start the NoOpReadLoop with mock", func() {
			go NoOpReadLoop(mc, ms)
			readyFail <- true
		})
		Convey("good start the NoOpReadLoop with mock", func() {
			go NoOpReadLoop(mc1, ms1)
		})
	})
}
