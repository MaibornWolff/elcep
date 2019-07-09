FROM library/golang:1.11.5-alpine AS build-env

LABEL version=1.11.5

RUN apk add --no-cache git gcc musl-dev