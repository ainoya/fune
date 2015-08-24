package actions

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/ainoya/fune/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/ainoya/fune/listener"
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
var RedisActionName = "redis"

// init function called once at launching program.
// In init() function, InstallAction is called and registers itself
// to `installedAction`
func init() {
	InstallAction(RedisActionName, &RedisAction{}, NewRedisAction)
}

// Name returns value `name` of struct `RedisAction`.
func (a *RedisAction) Name() string {
	return a.name
}

// Ch returns value `ch` of struct `RedisAction`.
func (a *RedisAction) Ch() chan interface{} {
	return a.ch
}

// On returns functions that register to redis
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
	fmt.Printf("redisAddr: %s\n", a.RedisAddr)
	go a.setAddressToRedis(redis)
}

// setAddressToRedis adds FQDN and IP:Port of container pair to Redis.
func (a *RedisAction) setAddressToRedis(redis *redis.Client) {
	for {
		e := <-a.out
		container, ipPort, err := repository.Listener.(*listener.DockerListener).ResolveIPPort(e)
		if err == nil {
			switch e.Status {
			case "start":
				redis.Set(fmt.Sprintf("%s.%s", container.Name[1:], a.BaseDomain), ipPort, 0)
				redis.Set(fmt.Sprintf("%s.%s", container.ID[0:7], a.BaseDomain), ipPort, 0)
				fmt.Printf("[redis] start container %s\n", a.RedisAddr)
			case "die":
				redis.Del(container.Name, container.ID[0:7])
				fmt.Printf("[redis] die container \n")
			}
		}
	}
}
