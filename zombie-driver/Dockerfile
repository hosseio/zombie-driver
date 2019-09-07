FROM golang:alpine AS builder
RUN apk update && apk add --no-cache alpine-sdk

ADD . /src
WORKDIR /src
RUN go mod download
RUN go mod verify
WORKDIR /src/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i -ldflags="-w -s"  -o ../bin/zombie-driver

FROM scratch
ENV PATH=/bin
COPY --from=builder /src/bin/zombie-driver /src/bin/zombie-driver
ENTRYPOINT ["/src/bin/zombie-driver"]