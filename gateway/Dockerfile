FROM golang:alpine AS build-env
RUN apk update && apk add --no-cache alpine-sdk
RUN apk --no-cache add ca-certificates
RUN adduser -D -g '' appuser

ADD . /src
WORKDIR /src
RUN go mod download
RUN go mod verify
WORKDIR /src/cmd
COPY ./cmd/config.yaml ../bin/config.yaml
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i -ldflags="-w -s"  -o ../bin/gateway

FROM scratch
ENV PATH=/bin
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /src/bin/gateway /src/bin/gateway
WORKDIR /src/bin
COPY ./cmd/config.yaml /src/bin/config.yaml
USER appuser
ENTRYPOINT ["/src/bin/gateway"]