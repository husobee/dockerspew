package content

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/unrolled/render.v1"
	"testing"
)

// TestGetContentTypeDecoder - get content type decoder
func TestGetContentTypeDecoder(t *testing.T) {
	Convey("test json decoder", t, func() {
		decoder := GetContentTypeDecoder(render.ContentJSON, &bytes.Buffer{})
		if _, ok := decoder.(*json.Decoder); !ok {
			t.Error("didn't get a json decoder")
		}
	})
	Convey("test default decoder", t, func() {
		decoder := GetContentTypeDecoder("", &bytes.Buffer{})
		if _, ok := decoder.(*json.Decoder); !ok {
			t.Error("didn't get a json decoder")
		}
	})
	Convey("test xml decoder", t, func() {
		decoder := GetContentTypeDecoder(render.ContentXML, &bytes.Buffer{})
		if _, ok := decoder.(*xml.Decoder); !ok {
			t.Error("didn't get a xml decoder")
		}
	})
}
