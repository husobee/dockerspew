package middlewares

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/pat"
	appctx "github.com/husobee/dockerspew/context"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHandler struct{}

func (th *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	So(context.Get(r, appctx.AcceptNegotiation), ShouldEqual, "text/xml")
	return
}

func TestContentNegotiateMiddleware(t *testing.T) {
	Convey("setup a negroni instance", t, func() {
		th := &TestHandler{}
		//setup a simple pat
		r := pat.New()
		r.Get("/test", th.ServeHTTP)

		// setup negroni
		n := negroni.Classic()
		n.Use(NewContentNegotiate(render.New()))

		Convey("setup a xml accept request", func() {
			req, err := http.NewRequest("GET", "/test", nil)
			req.Header.Add("Accept", "text/xml")
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()
			// attach router to negroni
			n.UseHandler(&TestHandler{})

			n.ServeHTTP(w, req)

		})
	})
}
