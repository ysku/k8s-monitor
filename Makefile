.PHONY: build
build:
	go build

.PHONY: init
init:
	go mod init

.PHONY: clean-test
clean-test:
	go clean -testcache

.PHONY: test
test: clean-test
	go test ./... -v

.PHONY: download
download:
	go mod download

.PHONY: verify
verify:
	go mod verify

.PHONY: lint
lint:
	golangci-lint run

.PHONY: push
push:
	docker build -t ysku/k8s-monitor:latest .
	docker push ysku/k8s-monitor
