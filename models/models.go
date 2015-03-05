// Package models - dockerspew models for accessing dockerapi
package models

import (
	"github.com/fsouza/go-dockerclient"
	"io"
)

var dockerNewClient = docker.NewClient

// NewDockerClient - Create a new global docker client based on a dockerEndpoint
func NewDockerClient(dockerEndpoint string) (*docker.Client, error) {
	return dockerNewClient(dockerEndpoint)
}

// DockerLogBuffer - implements io.Writer
type DockerLogBuffer struct {
	ContainerName []string
	ContainerID   string
	logChan       chan DockerLog
}

// DockerLog - a structure for log responses
type DockerLog struct {
	ContainerName []string `json:"name" xml:"Name"`
	ContainerID   string   `json:"id" xml:"ID"`
	LogMessage    string   `json:"log_message" xml:"LogMessage"`
}

//NewDockerLogBuffer - Create a new DockerLogBuffer with a channel of strings
func NewDockerLogBuffer(logChan chan DockerLog, containerName []string, containerID string) *DockerLogBuffer {
	return &DockerLogBuffer{
		ContainerName: containerName,
		ContainerID:   containerID,
		logChan:       logChan,
	}
}

// Write - implement io.Writer
func (dlb *DockerLogBuffer) Write(data []byte) (int, error) {
	dataLen := len(data)
	if dataLen == 0 {
		return 0, nil
	}

	dlb.logChan <- DockerLog{
		ContainerID:   dlb.ContainerID,
		ContainerName: dlb.ContainerName,
		LogMessage:    string(data),
	}
	return dataLen, nil
}

// GetLogChan - Get the log Channel
func (dlb *DockerLogBuffer) GetLogChan() chan DockerLog {
	return dlb.logChan
}

type DockerClientInterface interface {
	Logs(docker.LogsOptions) error
	ListContainers(docker.ListContainersOptions) ([]docker.APIContainers, error)
}

type MockDockerClient struct {
	MockLogs           func(docker.LogsOptions) error
	MockListContainers func(docker.ListContainersOptions) ([]docker.APIContainers, error)
}

func (mdc *MockDockerClient) Logs(opts docker.LogsOptions) error {
	if mdc.MockLogs != nil {
		return mdc.MockLogs(opts)
	}
	return nil
}

func (mdc *MockDockerClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	if mdc.MockListContainers != nil {
		return mdc.MockListContainers(opts)
	}
	return []docker.APIContainers{}, nil
}

//StreamContainerLogsById - Write out the container logs to out
func StreamContainerLogsById(dockerClient DockerClientInterface, id string, out io.Writer, stop chan bool) error {
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
	// need to rewrite dockerclient to accept a stop signal, i assume this will stream forever
	// not cool TODO
	return dockerClient.Logs(logOptions)
}

// GetAllContainers - Get all Containers in Docker
func GetAllContainers(dockerClient DockerClientInterface) ([]docker.APIContainers, error) {
	listContainersOptions := docker.ListContainersOptions{
		All: true,
	}
	return dockerClient.ListContainers(listContainersOptions)
}
