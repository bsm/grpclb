SRC=$(shell find . -name '*.go' -not -path '*vendor*')
PKG=$(shell glide nv)
PROTO=$(patsubst %.proto,%.pb.go,$(wildcard */*.proto))
VERSION=v0.3.1

TARGET_PKG=$(patsubst cmd/%/main.go,cmd/%,$(wildcard cmd/*/main.go))
TARGET_OS=linux darwin
TARGET_ARCH=amd64 386

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)

proto: $(PROTO)

touch-proto:
	touch $(wildcard */*.proto)

force-proto: touch-proto proto

.PHONY: default test vet all proto

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<

build-all:
	mkdir -p bin
	for pkg in $(TARGET_PKG); do \
		for os in $(TARGET_OS); do \
			for arch in $(TARGET_ARCH); do \
				env GOOS=$$os GOARCH=$$arch go build -o bin/$$(basename $$pkg)-$(VERSION)-$$os-$$arch ./$$pkg; \
			done; \
		done; \
	done
