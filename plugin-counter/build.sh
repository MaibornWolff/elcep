#!/bin/sh
ls
apk add --no-cache git gcc musl-dev
go get
go build --buildmode=plugin -o plugin-total.so