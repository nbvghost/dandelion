#!/bin/bash
### BEGIN INIT INFO
### END INIT INFO
echo 'build app'
export GOPATH='C:\Users\work\Documents\DandelionWork;C:\Users\work\go'
GOOS=linux go build -ldflags '-s -w'
go build -ldflags '-s -w'
