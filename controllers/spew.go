package controllers

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/websocket"
	"github.com/husobee/dockerspew/content"
	"github.com/husobee/dockerspew/models"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
	"strings"
)

// SpewController - Base Spew Controller
type SpewController struct {
	*Controller
	DockerClient      *docker.Client
	WebSocketUpgrader *websocket.Upgrader
	spew              chan []byte
}

// NewSpewController - Create a new Controller object
func NewSpewController(r *render.Render, dockerClient *docker.Client, webSocketUpgrader *websocket.Upgrader) *SpewController {
	log.Print("[DEBUG] Instantiation of a SpewController")
	return &SpewController{
		Controller:   NewController(r, webSocketUpgrader),
		DockerClient: dockerClient,
		spew:         make(chan []byte),
	}
}

// SpewHandler spews events from docker
func (sc *SpewController) SpewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[INFO] Starting SpewHandler")
	// get the contains variable from the URL, this will serve as a container name filter
	contains := r.URL.Query().Get("contains")
	log.Print("[DEBUG] SpewHandler - Pinging Docker Endpoint")
	err := sc.DockerClient.Ping()
	if err != nil {
		// error pinging docker endpoint, fail.
		log.Print("[DEBUG] Failed to Ping Docker Endpoint, err=", err.Error())
		// bail with 500
		sc.Respond(w, r, 500, content.NewBaseResponse("failure", "failed to ping", content.FailedDockerPingCode))
		return
	}
	// get all docker containers
	containers, err := models.GetAllContainers(sc.DockerClient)
	if err != nil {
		// error pinging docker endpoint, fail.
		log.Print("[DEBUG] Failed to get containers from Docker Endpoint, err=", err.Error())
		// bail with 500
		sc.Respond(w, r, 500, content.NewBaseResponse("failure", "failed to list containers", content.FailedDockerListContainersCode))
		return
	}
	var containerIDList []string
	for _, container := range containers {
		if !strings.Contains(container.Status, "Exit") {
			for _, name := range container.Names {
				if strings.Contains(name, contains) {
					containerIDList = append(containerIDList, container.ID)
				}
			}
		}
	}
	// Upgrade Connection to a websocket, coroutine out to websocket handler
	if wsConn, err := sc.UpgradeWebsocket(w, r); err == nil {
		webSocketConn := content.NewWebSocketConn(wsConn)
		go sc.WebSocketSpewHandler(webSocketConn, containerIDList...)
		return
	}
	sc.Respond(w, r, 500, content.NewBaseResponse("failure", "Failed to upgrade to websocket", content.FailedWebsocketUpgradeCode))
	return
}

// WebSocketSpewHandler - Handles upgraded websocket communications with client
func (sc *SpewController) WebSocketSpewHandler(wsConn *content.WebSocketConn, containerIds ...string) {
	// this is the log chan, where the DockerLogBuffer will spew logs
	var buf = models.NewDockerLogBuffer(make(chan string, 1024))
	var err error
	// throw away anything sent from client
	go content.NoOpReadLoop(wsConn)
	go func() {
		for {
			select {
			case message := <-buf.GetLogChan():
				log.Printf("[DEBUG] - message is %s\n", message)
				if writer, err := wsConn.NextWriter(websocket.TextMessage); err == nil {
					log.Println("[DEBUG] - got writer about to copy")
					if _, err = writer.Write([]byte(message)); err != nil {
						log.Println("[ERROR] Failed to copy Message to Websocket, err=", err)
					}
					log.Println("[DEBUG] Writing Message to Websocket, message=", message, " length=", len(message))
					log.Println("[DEBUG] - done copying, closing writer")
					if err = writer.Close(); err != nil {
						log.Println("[ERROR] Failed to close Websocket, err=", err)

					}
				} else {
					log.Println("[ERROR] Failed to Write Message to Websocket")
					return
				}
			}
		}
	}()
	// kick off streaming of container logs
	for _, containerID := range containerIds {
		go func(cID string) {
			err = models.StreamContainerLogsById(sc.DockerClient, cID, buf)
		}(containerID)
	}
	if err != nil {
		log.Println("[ERROR] Failed to Stream Container Logs, err=", err.Error())
		wsConn.Close()
		return
	}
}
