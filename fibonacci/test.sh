#!/bin/sh
cd "$GOPATH/src/fibonacci"

# Izvrši test
exec go test -v | go2xunit
