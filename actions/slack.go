package actions

import (
	"encoding/json"
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/ainoya/fune/listener"
	"net/http"
	"net/url"
)

// SlackAction is action which outputs the events that are received from listener.
type SlackAction struct {
	ch           chan interface{}
	out          chan *docker.APIEvents
	name         string
	BaseDomain   string
	IncomingURL  string `config:"slack-url" description:"Slack incoming webhook url"`
	dockerClient *docker.Client
}

// SlackActionName is used for identify name of itself.
var SlackActionName = "slack"

// init function called once at launching program.
// In init() function, InstallAction is called and registers itself
// to `installedAction`
func init() {
	InstallAction(SlackActionName, &SlackAction{}, NewSlackAction)
}

// Name returns value `name` of struct `SlackAction`.
func (a *SlackAction) Name() string {
	return a.name
}

// Ch returns value `ch` of struct `SlackAction`.
func (a *SlackAction) Ch() chan interface{} {
	return a.ch
}

// On returns functions that register to redis
func (a *SlackAction) On() func(event interface{}) {

	f := func(e interface{}) {
		a.out <- e.(*docker.APIEvents)
	}

	return f
}

//NewSlackAction returns instantiated `StdOutAction`.
func NewSlackAction() Action {
	a := &SlackAction{
		name:       SlackActionName,
		ch:         make(chan interface{}),
		out:        make(chan *docker.APIEvents),
		BaseDomain: "fune.dev", // TODO : make configurable from command line
	}

	return a
}

// Prepare is called when `actions` loads action instance.
// printMsg is called with goroutine to receive docker events
// from channel `SlackAction.ch` and print out them.
func (a *SlackAction) Prepare() {
	go a.postToSlack()
}

// setAddressToSlack adds FQDN and IP:Port of container pair to Slack.
func (a *SlackAction) postToSlack() {
	for {
		e := <-a.out
		a.post(e)
	}
}

// post message to slack incoming URL
// Reference: https://github.com/walter-cd/walter/blob/master/messengers/slack.go
// TODO : error handling
// TODO : add tests
func (a *SlackAction) post(e *docker.APIEvents) {

	container, _, _ := repository.Listener.(*listener.DockerListener).ResolveIPPort(e)

	var message string

	switch e.Status {
	case "start":
		message = fmt.Sprintf("[CREATED]container http://%s.%s is created.", container.Name[1:], a.BaseDomain)
	case "die":
		message = fmt.Sprintf("[DIED]container http://%s.%s is down.", container.Name[1:], a.BaseDomain)
	}

	params, _ := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: message,
	})

	resp, _ := http.PostForm(
		a.IncomingURL,
		url.Values{"payload": {string(params)}},
	)

	defer resp.Body.Close()
}
