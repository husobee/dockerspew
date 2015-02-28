package controllers

import (
	"github.com/husobee/dockerspew/content"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

// SpewController - Base Spew Controller
type SpewController struct {
	*Controller
}

// NewSpewController - Create a new Controller object
func NewSpewController(r *render.Render) *SpewController {
	log.Print("[DEBUG] Instantiation of a SpewController")
	return &SpewController{
		NewController(r),
	}
}

// SpewHandler spews events from docker
func (sc *SpewController) SpewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[INFO] Starting SpewHandler")
	sc.Respond(w, r, 200, content.NewBaseResponse("success", "success", 0))
	return
}
