# go-socket.io

[![GoDoc](http://godoc.org/github.com/gonetlib/go-socket.io?status.svg)](http://godoc.org/github.com/gonetlib/go-socket.io)
[![Build Status](https://github.com/gonetlib/go-socket.io/actions/workflows/ci.yaml/badge.svg)](https://github.com/gonetlib/go-socket.io/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonetlib/go-socket.io)](https://goreportcard.com/report/github.com/gonetlib/go-socket.io)

go-socket.io is library an implementation of [Socket.IO](http://socket.io) in Golang, which is a realtime application framework.

## JS Version compatible status

| JS Version | go-socket.io support |
| ---------- | -------------------- |
| 0.x        | :x:                  |
| 1.x        | :heavy_check_mark:   |
| 2.x        | :heavy_check_mark:   |
| 3.x        | :x:                  |
| 4.x        | :x:                  |

## Install

Install the package with:

```bash
go get github.com/gonetlib/socket.io
```

Import it with:

```go
import "github.com/gonetlib/socket.io"
```

and use `socketio` as the package name inside the code.

## Example

Please check more examples into folder in project for details. [Examples](./example/server/main.go)
