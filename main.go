// Package main - Server entrypoint for signature service
package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/hashicorp/logutils"
	"github.com/husobee/dockerspew/controllers"
	"github.com/husobee/dockerspew/middlewares"
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
	// get vars from viper env binding
	viper.SetEnvPrefix("dockerspew") // will be uppercased automatically
	viper.BindEnv("server_host")
	viper.BindEnv("log_level")
	// setup logging
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"},
		MinLevel: logutils.LogLevel(strings.ToUpper(viper.GetString("log_level"))),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

func main() {
	log.Print("[DEBUG] Starting Server, config options: server_host=", viper.GetString("server_host"), " log_level=", viper.GetString("log_level"))
	// setup renderer
	rend := render.New()
	// define routes
	r := pat.New()
	spewController := controllers.SpewController(rend)
	r.Get("/spew", spewController.Spew())
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
