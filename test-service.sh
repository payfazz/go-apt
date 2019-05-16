#!/bin/sh

#set -eu

go test -cover -coverprofile=c.out ./pkg/fazzcloud/*.go
go tool cover -html=c.out -o coverage.html
