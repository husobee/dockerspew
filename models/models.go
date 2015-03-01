// Package models - dockerspew models for accessing dockerapi
package models

import (
	"github.com/fsouza/go-dockerclient"
	"io"
)

// NewDockerClient - Create a new global docker client based on a dockerEndpoint
func NewDockerClient(dockerEndpoint string) (*docker.Client, error) {
	return docker.NewClient(dockerEndpoint)
}

// DockerLogBuffer - implements io.Writer
type DockerLogBuffer struct {
	logChan chan string
}

//NewDockerLogBuffer - Create a new DockerLogBuffer with a channel of strings
func NewDockerLogBuffer(logChan chan string) *DockerLogBuffer {
	return &DockerLogBuffer{
		logChan: logChan,
	}
}

// Write - implement io.Writer
func (dlb *DockerLogBuffer) Write(data []byte) (int, error) {
	dataLen := len(data)
	if dataLen == 0 {
		return 0, nil
	}
	dlb.logChan <- string(data)
	return dataLen, nil
}

// GetLogChan - Get the log Channel
func (dlb *DockerLogBuffer) GetLogChan() chan string {
	return dlb.logChan
}

//StreamContainerLogsById - Write out the container logs to out
func StreamContainerLogsById(dockerClient *docker.Client, id string, out io.Writer) error {
	logOptions := docker.LogsOptions{
		Container:    id,
		OutputStream: out,
		ErrorStream:  out,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   true,
		Tail:         "1",
		RawTerminal:  false,
	}
	return dockerClient.Logs(logOptions)
}

// GetAllContainers - Get all Containers in Docker
func GetAllContainers(dockerClient *docker.Client) ([]docker.APIContainers, error) {
	listContainersOptions := docker.ListContainersOptions{
		All: true,
	}
	return dockerClient.ListContainers(listContainersOptions)
}
