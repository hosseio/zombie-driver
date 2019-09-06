FROM golang:alpine AS builder
RUN apk update && apk add --no-cache alpine-sdk

ADD . /src
WORKDIR /src/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i -ldflags="-w -s"  -o ../bin/driver-location

FROM scratch
ENV PATH=/bin
COPY --from=builder /src/bin/driver-location /src/bin/driver-location
ENTRYPOINT ["/src/bin/driver-location"]