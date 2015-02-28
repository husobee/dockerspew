// Package controllers - Dockerspew controller package.  Controllers are structs that hold static configurations
// as well as handler functions that define the webservice.
package controllers

import (
	"bytes"
	"github.com/gorilla/context"
	"github.com/husobee/dockerspew/content"
	appctx "github.com/husobee/dockerspew/context"
	"gopkg.in/unrolled/render.v1"
	"io"
	"log"
	"net/http"
)

// Controller - Base Controller type, holds a common content responder, which is a renderer
type Controller struct {
	r content.Responder
}

// NewController - Create a new Controller object with a renderer
func NewController(r *render.Render) *Controller {
	log.Print("[DEBUG] Instanciating a new Controller; render=", r)
	return &Controller{
		r: content.NewResponder(r),
	}
}

// Respond - Respond to the request with the render's Respond method.
func (c *Controller) Respond(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	log.Print("[DEBUG] Response to Request; request=", r, "; code=", code, "; response=", v)
	c.r.Respond(w, r, code, v)
}

// Decode - Decode the request body based on the content-type, return an error, and the raw body in []byte
func (c *Controller) Decode(r *http.Request, v interface{}) ([]byte, error) {
	contentType := context.Get(r, appctx.ContentTypeNegotiation)
	// multiwriting the request body into two buffers, one for decoding and the other for logging, and returning
	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}
	mw := io.MultiWriter(b1, b2)
	io.Copy(mw, r.Body)
	requestBody := b2.Bytes()
	// logging out the request body
	log.Print("[DEBUG] Decoding Request Body Request; content_type=", contentType, "; request body=", string(requestBody))
	// creating a decoder with the first buffer
	decoder := content.GetContentTypeDecoder(contentType, b1)
	// returning the decoder.decode and the actual requestBody in []byte
	return requestBody, decoder.Decode(v)
}