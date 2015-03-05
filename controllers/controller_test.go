package controllers

import (
	"bytes"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"testing"
)

func TestControllerDecode(t *testing.T) {
	Convey("make a new controller", t, func() {
		rend := render.New()
		c := NewController(rend, &websocket.Upgrader{})
		Convey("do a decode on a json request", func() {
			req, err := http.NewRequest("GET", "/test", bytes.NewBuffer([]byte(`"test"`)))
			req.Header.Add("ContentType", "application/json")
			if err != nil {
				t.Error(err)
			}
			var v string
			c.Decode(req, &v)
			So(v, ShouldEqual, "test")
		})
	})
}
