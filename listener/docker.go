package listener

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"os"
)

type DockerListener struct {
	client  *docker.Client
	events  chan interface{}
	stopped chan struct{}
}

func NewDockerListener(c *docker.Client) *DockerListener {
	l := &DockerListener{}
	l.client = c

	return l
}

func GetDefaultClient() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	return client
}

func GetTLSClient() *docker.Client {
	endpoint := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
	return client
}

func (l *DockerListener) StartListen() {
	listener := make(chan interface{})
	l.events = listener
}

func (l *DockerListener) Events() <-chan interface{} {
	return l.events
}

func (l *DockerListener) Stopped() chan struct{} {
	return l.stopped
}

func (l *DockerListener) Stop() {
	close(l.stopped)
}
