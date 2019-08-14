#!/bin/bash
### BEGIN INIT INFO
### END INIT INFO
echo 'build app'
export GOPATH='C:\Users\work\Documents\DandelionWork;C:\Users\work\go'
GOARCH=386 go build
