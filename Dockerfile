FROM maibornwolff/elcep:builder-1.11.5 AS build-env

LABEL version=1.1

RUN mkdir -p /go/src/github.com/MaibornWolff/elcep
RUN mkdir -p /app

# download dependencies
ENV GO111MODULE=on
WORKDIR /go/src/github.com/MaibornWolff/elcep
COPY go.mod .
COPY go.sum .
RUN go mod download

# build elcep
COPY main /go/src/github.com/MaibornWolff/elcep/main
WORKDIR /go/src/github.com/MaibornWolff/elcep/main
RUN go test -v ./...
RUN go build -o /app/elcep

# build plugins
COPY plugins /go/src/github.com/MaibornWolff/elcep/plugins
WORKDIR /go/src/github.com/MaibornWolff/elcep/plugins
RUN for dir in */; do                                           \
        cd $dir;                                                \
        go test -v ./...;                                       \
        go build --buildmode=plugin -o /app/plugins/${dir%?}.so;\
        cd ..;                                                  \
    done

FROM alpine

WORKDIR /app
COPY --from=build-env /app/elcep /app/
COPY --from=build-env /app/plugins/*.so /app/plugins/

ENTRYPOINT ["./elcep"]
