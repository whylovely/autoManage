.PHONY: dev build test fmt vet check

dev:
	wails dev

build:
	wails build

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

check:
	go fmt ./...
	go vet ./...
	go test ./...
