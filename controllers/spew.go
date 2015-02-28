package controllers

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/husobee/dockerspew/content"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

// SpewController - Base Spew Controller
type SpewController struct {
	*Controller
	DockerClient *docker.Client
}

// NewSpewController - Create a new Controller object
func NewSpewController(r *render.Render, dockerClient *docker.Client) *SpewController {
	log.Print("[DEBUG] Instantiation of a SpewController")
	return &SpewController{
		Controller:   NewController(r),
		DockerClient: dockerClient,
	}
}

// SpewHandler spews events from docker
func (sc *SpewController) SpewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[INFO] Starting SpewHandler")
	err := sc.DockerClient.Ping()
	if err != nil {
		log.Print("[DEBUG] err=", err.Error())
		sc.Respond(w, r, 500, content.NewBaseResponse("failure", "failed to ping", content.FailedDockerPingCode))
		return
	}
	sc.Respond(w, r, 200, content.NewBaseResponse("success", "success", content.SuccessCode))
	return
}
