#!/bin/bash

gofmt -s -w ./apiserver
go tool fix ./apiserver
go tool vet ./apiserver

CGO_ENABLED=0 go test ./apiserver

CGO_ENABLED=0 go install ./apiserver
