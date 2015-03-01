// Package main - Server entrypoint for signature service
package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/hashicorp/logutils"
	"github.com/husobee/dockerspew/content"
	"github.com/husobee/dockerspew/controllers"
	"github.com/husobee/dockerspew/middlewares"
	"github.com/husobee/dockerspew/models"
	"github.com/spf13/viper"
	"gopkg.in/unrolled/render.v1"
	"log"
	"os"
	"strings"
)

func init() {
	// set some defaults
	viper.SetDefault("server_host", ":8080")
	viper.SetDefault("log_level", "WARN")
	viper.SetDefault("docker_endpoint", "unix:///var/run/docker.sock")
	// get vars from viper env binding
	viper.SetEnvPrefix("dockerspew") // will be uppercased automatically
	viper.BindEnv("server_host")
	viper.BindEnv("log_level")
	viper.BindEnv("docker_endpoint")
	// setup logging
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"},
		MinLevel: logutils.LogLevel(strings.ToUpper(viper.GetString("log_level"))),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

func main() {
	log.Print("[DEBUG] Starting Server, config options: docker_endpoint:",
		viper.GetString("docker_endpoint"),
		" server_host=",
		viper.GetString("server_host"),
		" log_level=",
		viper.GetString("log_level"))
	// setup renderer
	rend := render.New()
	// setup a docker client
	dockerClient, err := models.NewDockerClient(viper.GetString("docker_endpoint"))
	if err != nil {
		log.Fatalln("[PANIC] Docker enpoint not accessable, bailing")
	}
	// define routes
	r := pat.New()
	// setup a websocketUpgrader
	webSocketUpgrader := content.NewWebSocketUpgrader(1024, 1024)
	// setup our spew controller
	spewController := controllers.NewSpewController(rend, dockerClient, webSocketUpgrader)
	r.Get("/spew", spewController.SpewHandler)
	r.Get("/spew/", spewController.SpewHandler)
	// startup classic negroni
	n := negroni.Classic()
	n.Use(middlewares.NewContentNegotiate(rend))
	// attach router to negroni
	n.UseHandler(r)
	// run negroni
	runServer(n)
}

//runServer - run the server, broken out for unit tests
var runServer = func(n *negroni.Negroni) {
	log.Print("[DEBUG] Server starting to accept requests")
	n.Run(viper.GetString("server_host"))
}
