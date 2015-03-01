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
	var containerList []docker.APIContainers
	for _, container := range containers {
		if !strings.Contains(container.Status, "Exit") {
			for _, name := range container.Names {
				if strings.Contains(name, contains) {
					containerList = append(containerList, container)
				}
			}
		}
	}
	// Upgrade Connection to a websocket, coroutine out to websocket handler
	if wsConn, err := sc.UpgradeWebsocket(w, r); err == nil {
		webSocketConn := content.NewWebSocketConn(wsConn)
		go sc.WebSocketSpewHandler(r, webSocketConn, containerList...)
		return
	}
	sc.Respond(w, r, 500, content.NewBaseResponse("failure", "Failed to upgrade to websocket", content.FailedWebsocketUpgradeCode))
	return
}

// WebSocketSpewHandler - Handles upgraded websocket communications with client
func (sc *SpewController) WebSocketSpewHandler(r *http.Request, wsConn *content.WebSocketConn, containers ...docker.APIContainers) {
	// this is the log chan, where the DockerLogBuffer will spew logs
	var logChan = make(chan models.DockerLog, 1024)
	var err error
	// throw away anything sent from client
	stopReading := make(chan bool)
	stopStreaming := make(map[string]chan bool)
	go content.NoOpReadLoop(wsConn, stopReading)
	go func() {
		for {
			select {
			// wait for the logChan
			case message := <-logChan:
				// on message get the next writer
				log.Printf("[DEBUG] - message is %s\n", message)
				if writer, err := wsConn.NextWriter(websocket.TextMessage); err == nil {
					log.Println("[DEBUG] - got writer about to copy")
					// write the response
					if _, err := sc.WSRespond(writer, r, message); err != nil {
						log.Println("[ERROR] Failed to copy Message to Websocket, err=", err)
					}
					log.Println("[DEBUG] - done copying, closing writer")
					// close the writer
					if err = writer.Close(); err != nil {
						log.Println("[ERROR] Failed to close Websocket, err=", err)
					}
				} else {
					// failed to write to the websocket, probably dead, end this goroutine, and tell read to stop too
					wsConn.KillChan <- true
				}
			case <-wsConn.KillChan:
				stopReading <- true
				for _, v := range stopStreaming {
					v <- true
					close(v)
				}
				close(logChan)
				log.Println("[ERROR] Cleaning up after websocket")
				return
			}
		}
	}()

	// kick off streaming of container logs
	for _, container := range containers {
		var buf = models.NewDockerLogBuffer(logChan, container.Names, container.ID)
		stopStreaming[container.ID] = make(chan bool)
		go func(cID string) {
			err = models.StreamContainerLogsById(sc.DockerClient, cID, buf, stopStreaming[cID])
		}(container.ID)
	}
	if err != nil {
		log.Println("[ERROR] Failed to Stream Container Logs, err=", err.Error())
		wsConn.Close()
		return
	}
}
