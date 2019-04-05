#!/bin/sh

#set -eu

go test -cover -coverprofile=c.out ./pkg/fazzcommon/formatter/*.go
go tool cover -html=c.out -o coverage.html
