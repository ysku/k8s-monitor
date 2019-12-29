.PHONY: build
build:
	go build

push:
	docker build -t ysku/k8s-monitor:latest .
	docker push ysku/k8s-monitor

