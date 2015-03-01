package content

import (
	"encoding/xml"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

const (
	// SuccessCode Response Status
	SuccessCode int = iota
	// FailedDockerPingCode Response Status
	FailedDockerPingCode
	// FailedWebsocketUpgradeCode Response Status
	FailedWebsocketUpgradeCode
)

// Responder - Overloaded unrolled/render.v1 Render
type Responder struct {
	*render.Render
}

// NewResponder - Create a new "responder" we can attach to controllers, so we can perform
// dynamic rendering of responses
func NewResponder(r *render.Render) Responder {
	return Responder{
		r,
	}
}

// Respond - Based on the request's Accept Header, return the proper rendering.
func (r Responder) Respond(w http.ResponseWriter, accept interface{}, s int, v interface{}) {
	log.Print("[DEBUG] Responder.Respond - Accept is ", accept)
	switch accept {
	case render.ContentJSON:
		log.Print("[DEBUG] Responder.Respond - rendering JSON of v=", v)
		r.JSON(w, s, v)
	case render.ContentXML:
		log.Print("[DEBUG] Responder.Respond - rendering XML of v=", v)
		r.XML(w, s, v)
	default:
		log.Print("[DEBUG] Responder.Respond - rendering JSON of v=", v)
		r.JSON(w, s, v)
	}
}

// BaseResponse - this is the basic structure of all responses, for consistancy
type BaseResponse struct {
	XMLName xml.Name `json:"-" xml:"Response"`
	Status  string   `json:"status" xml:"Status"`
	Message string   `json:"message" xml:"Message"`
	Code    int      `json:"code" xml:"Code"`
}

// NewBaseResponse - create a new base response
func NewBaseResponse(status, message string, code int) BaseResponse {
	return BaseResponse{
		Status:  status,
		Message: message,
		Code:    code,
	}
}
