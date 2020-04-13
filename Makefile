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
	docker build -t ysku/my-k8s-custom-controller:latest .
	docker push ysku/my-k8s-custom-controller

.PHONY: apply
apply:
	@kubectl apply -f deploy/

.PHONY: delete
delete:
	-@kubectl delete -f deploy/

.PHONY: logs
logs:
	@kubectl get pod --selector=app=my-k8s-custom-controller --output=jsonpath={.items..metadata.name} | xargs kubectl logs -f
