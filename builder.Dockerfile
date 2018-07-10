FROM library/golang:1.10.2-alpine AS build-env

LABEL version=1.20.2

RUN apk add --no-cache git gcc musl-dev