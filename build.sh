#!/bin/zsh

export GOOS=windows
export GOARCH=amd64
go build -ldflags "-s -w"  -o cpss.exe main.go

export GOOS=windows
export GOARCH=amd64
go build -ldflags "-s -w"  -o cpss.exe main.go