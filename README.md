Fune
==========

[![wercker status](https://app.wercker.com/status/ffc850729c83561933ed7651a25634b7/m "wercker status")](https://app.wercker.com/project/bykey/ffc850729c83561933ed7651a25634b7)

An Event emitting agent triggered by docker daemon

Overview
----------

![](https://i.gyazo.com/d242e18e636b06add0b9ebc769188232.png)

Fune is an agent emitting some actions when docker daemon causes events.


Implemented actions are for example:

- **TODO** Register/Deregister container's information to `coreos/etcd`
- **Experimental** Register/Deregister container's addresses to `openresty` to make containers accessible via proxy

You can check more actions in "Actions list" section.

How to run
-----------

To run fune agent, you need to run `fune-agent` in your docker hosts.

### Setup

First, get binary from [Releases · ainoya/fune](https://github.com/ainoya/fune/releases) .

```
wget https://github.com/ainoya/fune/releases/latest/v0.1.0/fune-agent \
  -O /usr/local/bin/fune-agent && chmod +x /usr/local/bin/fune-agent
```

### Run

You can enable an action with `--actions` option when you runs `fune-agent`.

For example, if you enable `stdout` action, `fune-agent` output docker events
into STDOUT simply.

```
$ ./build && ./bin/fune-agent --actions=stdout
$  docker run -it ubuntu echo 'hello'
$ (fune-agent output)
Status: create ID: de88fcafffb616ceb1e9191162affe55174f3e0531c0ad9fffc83d535313aeb6 From: ubuntu Time: 1439734180
Status: start ID: de88fcafffb616ceb1e9191162affe55174f3e0531c0ad9fffc83d535313aeb6 From: ubuntu Time: 1439734180
Status: die ID: de88fcafffb616ceb1e9191162affe55174f3e0531c0ad9fffc83d535313aeb6 From: ubuntu Time: 1439734180
```

Build
--------

```
go get github.com/tools/godep
go get github.com/ainoya/fune
cd ${GOPATH}/src/github.com/ainoya/fune
godep restore
./build # run a shell script to build binary
```

Actions list
--------------

### Slack

`fune-agent` supports slack notification.

```
./bin/fune-agent -actions=slack --slack-url="Incoming URL"
```

Then, `fune-agent` notifies docker events to slack channel like below.

![](https://i.gyazo.com/718b803555bb980545866b26e6a3622f.png)

### Redis

It also can register container information to Redis as FQDN and container
ip/port pair. With using these container informations stored in Redis, you can implement redis-backed dynamic proxy (like openresty) easily.

```
# Run Openresty container
$ docker run --name openresty -i -t -d -p 80:80 -p 6379:6379 quay.io/ainoya/openresty-dynamic-upstream
$ ./bin/fune-agent --actions=redis --redis-addr=0.0.0.0:6379 --base-domain=fune.dev
# Run container
$  docker run -it ubuntu -p 80:8000 sh -c "python -m SimpleHTTPserver"
#
$ redis-cli
# There are stored FQDN corresponds to container addresses.
127.0.0.1:6379> keys *
1) "high_stallman.fune.dev"
2) "27f180f.fune.dev"
# You can access container with FQDN stored in Redis.
$ curl http://high_stallman.fune.dev
200 OK...
```

Test
------

You can run unit/integration tests with `./test` written in shell script:

```
./test
```

TODO
-----

- [ ] Make actions pluggable


Reference
------------

- [coreos/etcd](https://github.com/coreos/etcd): An implementation of Fune is heavily referencing etcd's.
