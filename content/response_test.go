package content

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/unrolled/render.v1"
	"net/http/httptest"
	"testing"
)

// TestRespond
func TestRespond(t *testing.T) {
	Convey("Create a new Responder", t, func() {
		rend := render.New()
		r := NewResponder(rend)
		w := httptest.NewRecorder()
		Convey("call Respond with json acccept type", func() {
			r.Respond(w, "application/json", 200, "hi")
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldEqual, "\"hi\"")
		})
		Convey("call Respond with random acccept type", func() {
			r.Respond(w, "application/j", 200, "hi")
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldEqual, "\"hi\"")
		})
		Convey("call Respond with xml acccept type", func() {
			r.Respond(w, "text/xml", 200, "hi")
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldEqual, "<string>hi</string>")
		})
	})
}

type testWriteCloser struct {
	Wrote []byte
}

func (twc *testWriteCloser) Write(b []byte) (int, error) {
	twc.Wrote = b
	return len(b), nil
}

func (twc *testWriteCloser) Close() error {
	return nil
}

// TestRespondWS
func TestRespondWS(t *testing.T) {
	Convey("Create a new Responder", t, func() {
		rend := render.New()
		r := NewResponder(rend)
		w := &testWriteCloser{}
		Convey("call Respond with json acccept type", func() {
			r.WSRespond(w, "application/json", "hi")
			So(string(w.Wrote), ShouldEqual, "\"hi\"")
		})
		Convey("call Respond with json acccept type, json error", func() {
			oldJsonMarshal := jsonMarshal
			defer func() { jsonMarshal = oldJsonMarshal }()
			jsonMarshal = func(v interface{}) ([]byte, error) {
				return []byte{}, errors.New("failure")
			}
			_, err := r.WSRespond(w, "application/json", "hi")
			So(err, ShouldNotBeNil)
		})
		Convey("call Respond with random acccept type", func() {
			r.WSRespond(w, "application/j", "hi")
			So(string(w.Wrote), ShouldEqual, "\"hi\"")
		})
		Convey("call Respond with xml acccept type", func() {
			r.WSRespond(w, "text/xml", "hi")
			So(string(w.Wrote), ShouldEqual, "<string>hi</string>")
		})
		Convey("call Respond with xml acccept type, xml error", func() {
			oldXMLMarshal := xmlMarshal
			defer func() { xmlMarshal = oldXMLMarshal }()
			xmlMarshal = func(v interface{}) ([]byte, error) {
				return []byte{}, errors.New("failure")
			}
			_, err := r.WSRespond(w, "text/xml", "hi")
			So(err, ShouldNotBeNil)
		})
	})
}

// TestNewBaseResponse - test structure of base response
func TestNewBaseRespone(t *testing.T) {
	Convey("Create a new Base Response", t, func() {
		r := NewBaseResponse("success", "message", 200)
		So(r.Code, ShouldEqual, 200)
		So(r.Status, ShouldEqual, "success")
		So(r.Message, ShouldEqual, "message")
	})
}
