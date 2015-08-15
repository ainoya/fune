package listener

import (
	"testing"
)

func TestNewDockerListenerWithDefault(t *testing.T) {
	c := GetDefaultClient()
	NewDockerListener(c)
}

func TestNewDockerListenerWithTLSClient(t *testing.T) {
	c := GetTLSClient()
	NewDockerListener(c)
}
