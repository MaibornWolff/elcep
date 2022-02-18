FROM library/golang:1.17.2-alpine3.14 AS build-env

LABEL version=1.2

RUN apk add --no-cache git gcc musl-dev

RUN mkdir -p /src/main
RUN mkdir -p /app
WORKDIR /src

COPY main/go.mod main/.
COPY main/plugin/go.mod main/plugin/
COPY main/config/go.mod main/config/

COPY main/go.sum .
COPY main/plugin/go.sum main/plugin/
COPY main/config/go.sum main/config/

WORKDIR /src/main

RUN go mod download

# build elcep
COPY main .

RUN go test -v ./... 
RUN go build -o /app/elcep

WORKDIR /src

# build plugins
COPY plugins plugins
WORKDIR /src/plugins
RUN for dir in */; do                                           \
        cd $dir;                                                \
        go test -v ./...;                                       \
        go build --buildmode=plugin -o /app/plugins/${dir%?}.so;\
        cd ..;                                                  \
    done

FROM alpine:3.14

# create non root user
RUN addgroup -S nonroot && \
adduser -S nonroot -G nonroot

WORKDIR /app
COPY --chown=nonroot:nonroot --from=build-env /app/elcep /app/
COPY --chown=nonroot:nonroot --from=build-env /app/plugins/*.so /app/plugins/

USER nonroot

ENTRYPOINT ["./elcep"]
