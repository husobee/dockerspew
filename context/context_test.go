package context

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConstants(t *testing.T) {
	Convey("ContentTypeNegotiation and AcceptNegotiation context keys should be set", t, func() {
		So(ContentTypeNegotiation, ShouldEqual, 0)
		So(AcceptNegotiation, ShouldEqual, 1)
	})
}
