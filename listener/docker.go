package listener

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"os"
)

// DockerListener implements listening docker events.
type DockerListener struct {
	client  *docker.Client
	events  chan interface{}
	stopped chan struct{}
}

// NewDockerListener returns instantiated `DockerListner`
func NewDockerListener(c *docker.Client) *DockerListener {
	l := &DockerListener{}
	l.client = c

	return l
}

// GetDefaultClient gets `docker.Client` which uses unix domain socket.
func GetDefaultClient() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	return client
}

// GetTLSClient gets `docker.Client` which uses TCP connection.
func GetTLSClient() *docker.Client {
	endpoint := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
	return client
}

// StartListen prepares to listen docker events from Docker API.
func (l *DockerListener) StartListen() {
	listener := make(chan *docker.APIEvents)
	l.events = make(chan interface{})
	l.client.AddEventListener(listener)

	go func() {
		for {
			e, ok := <-listener

			if !ok {
				close(l.events)
				return
			}
			l.events <- e
		}
	}()
}

// Events returns a channel of docker events.
func (l *DockerListener) Events() <-chan interface{} {
	return l.events
}

// Stopped returns a channel to check if it is closed
func (l *DockerListener) Stopped() chan struct{} {
	return l.stopped
}

// Stop notifiies that DockerListner is closed as closing `l.stopped` channel.
func (l *DockerListener) Stop() {
	close(l.stopped)
}
