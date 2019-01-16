FROM maibornwolff/elcep:builder-1.10.2 AS build-env

LABEL version=0.5

RUN mkdir -p /go/src/github.com/MaibornWolff/elcep
RUN mkdir -p /go/src/github.com/MaibornWolff/elcep/plugins

COPY . /go/src/github.com/MaibornWolff/elcep
WORKDIR /go/src/github.com/MaibornWolff/elcep

# build elcep
RUN go get -d -v -t ./...
RUN go test -v ./...
RUN go build -o elcep

# build plugins
WORKDIR /go/src/github.com/MaibornWolff/elcep/plugins
RUN for dir in */; do OUTPUT_DIR="$(pwd)" ./${dir}build.sh; done

FROM alpine

WORKDIR /app
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/elcep /app/
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/plugins /app/plugins
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/conf /app/conf

ENTRYPOINT ["./elcep"]
CMD [""]