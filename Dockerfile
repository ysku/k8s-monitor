FROM golang:1.13.4-alpine AS build-env

RUN apk update && apk --no-cache add curl git ca-certificates alpine-sdk

ADD . /go/src/github.com/ysku/my-k8s-custom-controller
WORKDIR /go/src/github.com/ysku/my-k8s-custom-controller

RUN make build GOARCH=amd64 CGO_ENABLED=0 GOOS=linux

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /apps/
COPY --from=build-env /go/src/github.com/ysku/my-k8s-custom-controller /apps/

ENTRYPOINT ["./my-k8s-custom-controller"]
