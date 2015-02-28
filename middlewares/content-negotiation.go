package middlewares

import (
	"github.com/gorilla/context"
	appctx "github.com/husobee/dockerspew/context"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

// ContentNegotiate - Handler Struct that implements negroni.Handler
type ContentNegotiate struct {
	r *render.Render
}

// NewContentNegotiate - Return a new *ContentNegotiate, initialization of middleware
func NewContentNegotiate(r *render.Render) *ContentNegotiate {
	log.Print("[DEBUG] Instanciating a content negotiation middleware")
	return &ContentNegotiate{
		r: r,
	}
}

// ServeHTTP - Implement negroni.Handler interface, set the context based on the request header.
func (cn *ContentNegotiate) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Print("[DEBUG] New request, pulling Content-Type and Accept from Headers and stashing in context")
	context.Set(r, appctx.ContentTypeNegotiation, r.Header.Get("Content-Type"))
	context.Set(r, appctx.AcceptNegotiation, r.Header.Get("Accept"))
	log.Print("[DEBUG] Calling next")
	next(rw, r)

}
