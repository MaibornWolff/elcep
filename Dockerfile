FROM maibornwolff/elcep:builder-1.10.2 AS build-env

LABEL version=0.5

RUN mkdir -p /go/src/github.com/MaibornWolff/elcep
RUN mkdir -p /go/src/github.com/MaibornWolff/elcep/plugins

# get some dependencies before copying the source
# allows caching those deps =)
RUN go get -v -d -t gopkg.in/alecthomas/kingpin.v2
RUN go get -v -d -t gopkg.in/go-playground/assert.v1
RUN go get -v -d -t gopkg.in/yaml.v2
RUN go get -v -d -t github.com/olivere/elastic
RUN go get -v -d -t github.com/prometheus/client_golang/prometheus
RUN go get -v -d -t github.com/golang/mock/gomock
RUN go get -v -d -t github.com/mitchellh/hashstructure

COPY . /go/src/github.com/MaibornWolff/elcep
WORKDIR /go/src/github.com/MaibornWolff/elcep

# build elcep
RUN go get -d -v -t ./...
RUN go test -v ./...
RUN go build -o elcep

# build shipped_plugins
WORKDIR /go/src/github.com/MaibornWolff/elcep/shipped_plugins
RUN for dir in */; do chmod +x ${dir}build.sh && OUTPUT_DIR="$(pwd)" ./${dir}build.sh; done

FROM alpine

WORKDIR /app
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/elcep /app/
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/shipped_plugins/*.so /app/plugins/

ENTRYPOINT ["./elcep"]
