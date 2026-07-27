.PHONY: dev build build-macos build-windows build-linux test test-race coverage frontend-check fmt vet check

COVERAGE_FILE ?= coverage.out

dev:
	wails dev

build:
	wails build

build-macos:
	wails build -platform darwin/universal

build-windows:
	wails build -platform windows/amd64 -nsis

build-linux:
	wails build -platform linux/amd64

test:
	go test ./...

test-race:
	go test -race ./...

coverage:
	go test -coverprofile=$(COVERAGE_FILE) ./internal/service ./internal/storage ./internal/scheduler
	go tool cover -func=$(COVERAGE_FILE)

frontend-check:
	cd frontend && npm run build

fmt:
	go fmt ./...

vet:
	go vet ./...

check:
	go fmt ./...
	go vet ./...
	go test ./...
	cd frontend && npm run build
