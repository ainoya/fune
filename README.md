Fune
==========

[![wercker status](https://app.wercker.com/status/ffc850729c83561933ed7651a25634b7/m "wercker status")](https://app.wercker.com/project/bykey/ffc850729c83561933ed7651a25634b7)

An Event emitting agent triggered by docker daemon

Overview
----------

Fune is an agent emitting some actions when docker daemon causes events.

Implemented actions are:

- Register/Deregister container's information to `coreos/etcd`
- Register/Deregister container's addresses to `hipache` to make containers accessible via proxy

How to run
-----------

To run fune agent, you need to run `fune-agent` in your docker hosts.

### Setup

TBD

### Run

TBD


Build
--------

Build requires [tools/godep](https://github.com/tools/godep).

```
go get github.com/tools/godep
go get github.com/ainoya/fune
cd ${GOPATH}/src/github.com/ainoya/fune
godep restore
./build # run a shell script to build binary
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
