#!/bin/zsh
export GOARCH=amd64

export GOOS=windows
go build -ldflags "-s -w"  -o bin/cpss.exe main.go

export GOOS=linux
go build -ldflags "-s -w"  -o bin/cpssux main.go

export GOOS=darwin
go build -ldflags "-s -w"  -o bin/cpss main.go