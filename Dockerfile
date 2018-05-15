FROM library/golang:alpine AS build-env

RUN apk add --no-cache git gcc musl-dev

RUN mkdir -p /go/src/github.com/MaibornWolff/elcep
RUN mkdir -p /go/src/github.com/MaibornWolff/elcep/plugins

COPY . /go/src/github.com/MaibornWolff/elcep
WORKDIR /go/src/github.com/MaibornWolff/elcep

RUN go get -d -v -t ./...
RUN go test -v ./...
RUN go build -o elcep

FROM alpine

WORKDIR /app
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep /app/
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/plugins /app/plugins
COPY --from=build-env /go/src/github.com/MaibornWolff/elcep/conf /app/conf

ENTRYPOINT ["./elcep"]
CMD [""]