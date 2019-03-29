#!/bin/sh

out_dir=${OUTPUT_DIR:-.}
cd -P -- $(dirname -- "$0")

go get -d -v -t ./...
go test -v ./...
go build --buildmode=plugin -o "$out_dir/bucket.so"
