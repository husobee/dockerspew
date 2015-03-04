package models

import (
	"errors"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

// TestNewDockerClient - to get 100% coverage
func TestNewDockerClient(t *testing.T) {
	Convey("make a new docker client", t, func() {
		oldDockerNewClient := dockerNewClient
		dockerNewClient = func(e string) (*docker.Client, error) {
			return &docker.Client{}, nil
		}
		defer func() { dockerNewClient = oldDockerNewClient }()

		dc, err := NewDockerClient("123")
		So(dc, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

// TestStreamContainerLogsById - test streaming container logs
func TestStreamContainerLogsById(t *testing.T) {
	Convey("mocked docker client to stream logs", t, func() {
		stopChan := make(chan bool)
		StreamContainerLogsById(&MockDockerClient{}, "123", os.Stdout, stopChan)
		close(stopChan)
	})
	Convey("mocked docker client to error ", t, func() {
		stopChan := make(chan bool)
		err := StreamContainerLogsById(&MockDockerClient{
			MockLogs: func(opts docker.LogsOptions) error {
				return errors.New("fail")
			},
		}, "123", os.Stdout, stopChan)
		So(err, ShouldNotBeNil)
		close(stopChan)
	})
}

// TestGetAllContainers - test getting containers from docker
func TestGetAllContainers(t *testing.T) {
	Convey("mocked docker client to GetContainers", t, func() {
		_, _ = GetAllContainers(&MockDockerClient{})
	})
	Convey("mocked docker client to return error", t, func() {
		_, err := GetAllContainers(&MockDockerClient{
			MockListContainers: func(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
				return []docker.APIContainers{}, errors.New("fail")
			},
		})
		So(err, ShouldNotBeNil)
	})
}

// TestDockerLogBufferWrite - test writing of a docker log buffer
func TestDockerLogBufferWrite(t *testing.T) {
	Convey("setup the logchan, and dockerlogbuffer", t, func() {
		logChan := make(chan DockerLog)
		dockerLogBuffer := NewDockerLogBuffer(logChan, []string{"test"}, "123")
		Convey("try writing to the docker log buffer", func() {
			go func() {
				_, err := dockerLogBuffer.Write([]byte("this is a test"))
				if err != nil {
					t.Error("failed to write to log buffer")
				}
			}()
			message := <-dockerLogBuffer.GetLogChan()
			So(message, ShouldNotBeNil)
			So(message.LogMessage, ShouldEqual, "this is a test")
			So(message.ContainerID, ShouldEqual, "123")
			So(message.ContainerName[0], ShouldEqual, "test")
		})
		Convey("try writing no data to the docker log buffer", func() {
			res, err := dockerLogBuffer.Write([]byte(""))
			if err != nil {
				t.Error("failed to write to log buffer")
			}
			So(res, ShouldEqual, 0)
		})
		close(logChan)
	})
}
