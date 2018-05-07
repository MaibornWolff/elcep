FROM library/golang:alpine AS build-env

RUN apk add --no-cache git gcc musl-dev

RUN mkdir -p /go/src/elcep/plugins
RUN mkdir /go/src/elcep/conf

COPY . /go/src/elcep/
WORKDIR /go/src/elcep

RUN go get -d -v -t ./...
RUN go test -v ./...
RUN go build -o elcep

FROM alpine

WORKDIR /app
COPY --from=build-env /go/src/elcep/elcep /app/
COPY --from=build-env /go/src/elcep/plugins /app/plugins
COPY --from=build-env /go/src/elcep/conf /app/conf

ENTRYPOINT ["./elcep"]
CMD [""]