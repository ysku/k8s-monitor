FROM golang:1.13.4-alpine AS build-env

RUN apk update && apk --no-cache add curl git ca-certificates alpine-sdk

ADD . /go/src/github.com/ysku/k8s-monitor
WORKDIR /go/src/github.com/ysku/k8s-monitor

RUN make build GOARCH=amd64 CGO_ENABLED=0 GOOS=linux

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /apps/
COPY --from=build-env /go/src/github.com/ysku/k8s-monitor /apps/

ENTRYPOINT ["./k8s-monitor"]
