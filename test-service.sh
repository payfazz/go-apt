#!/bin/sh

#set -eu

go test -cover -coverprofile=c.out ./pkg/fazzcommon/encryption/*.go
go tool cover -html=c.out -o coverage.html
