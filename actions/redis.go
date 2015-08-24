package actions

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/ainoya/fune/listener"
	"gopkg.in/redis.v3"
)

// RedisAction is action which outputs the events that are received from listener.
type RedisAction struct {
	ch           chan interface{}
	out          chan *docker.APIEvents
	name         string
	BaseDomain   string `config:"redis-base-domain" description:"Base domain name: ex) fune.dev"`
	RedisAddr    string `config:"redis-addr" description:"Redis address: ex) 0.0.0.0:6379"`
	dockerClient *docker.Client
}

// RedisActionName is used for identify name of itself.
var RedisActionName = "stdout"

// init function called once at launching program.
// In init() function, InstallAction is called and registers itself
// to `installedAction`
func init() {
	InstallAction(RedisActionName, &StdOutAction{}, NewStdOutAction)
}

// Name returns value `name` of struct `RedisAction`.
func (a *RedisAction) Name() string {
	return a.name
}

// Ch returns value `ch` of struct `RedisAction`.
func (a *RedisAction) Ch() chan interface{} {
	return a.ch
}

// On returns functions that outputs to STDOUT.
func (a *RedisAction) On() func(event interface{}) {

	f := func(e interface{}) {
		a.out <- e.(*docker.APIEvents)
	}

	return f
}

//NewRedisAction returns instantiated `StdOutAction`.
func NewRedisAction() Action {
	a := &RedisAction{
		name:       RedisActionName,
		ch:         make(chan interface{}),
		out:        make(chan *docker.APIEvents),
		BaseDomain: "fune.dev", // TODO : make configurable from command line
	}

	return a
}

// Prepare is called when `actions` loads action instance.
// printMsg is called with goroutine to receive docker events
// from channel `RedisAction.ch` and print out them.
func (a *RedisAction) Prepare() {
	redis := redis.NewClient(
		&redis.Options{
			Addr:     a.RedisAddr,
			Password: "",
			DB:       0,
		},
	)
	go a.setAddressToRedis(redis)
}

func (a *RedisAction) setAddressToRedis(redis *redis.Client) {
	for {
		e := <-a.out
		container, ipPort, err := repository.Listener.(*listener.DockerListener).ResolveIPPort(e)
		if err == nil {
			switch e.Status {
			case "create":
				redis.Set(container.Name, ipPort, 0)
				redis.Set(container.ID[0:7], ipPort, 0)
			case "destroy":
				redis.Del(container.Name, container.ID[0:7])
			}
		}
	}
}
