package main

import (
	"errors"
	"github.com/codegangsta/negroni"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMain(t *testing.T) {
	Convey("change runserver to not actually run the server", t, func() {
		oldRunMain := runServer
		runServer = func(n *negroni.Negroni) {}
		defer func() { runServer = oldRunMain }()
		Convey("run main for coverage", func() {
			main()
		})
		Convey("run main failed dockerclient", func() {
			oldNewDockerClient := modelsNewDockerClient
			modelsNewDockerClient = func(s string) (*docker.Client, error) {
				return nil, errors.New("fail")
			}
			defer func() { modelsNewDockerClient = oldNewDockerClient }()
			defer func() {
				if r := recover(); r != nil {
					t.Log("Recovered in f", r)
				}
			}()
			main()
		})
	})
}
