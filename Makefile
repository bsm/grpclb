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

bin/%.gz: bin/%
	gzip -q -c $< > $@

package-all: \
	bin/grpc-lb-consul-$(VERSION)-linux-amd64.gz \
	bin/grpc-lb-consul-$(VERSION)-darwin-amd64.gz \
	bin/grpc-lb-client-$(VERSION)-linux-amd64.gz \
	bin/grpc-lb-client-$(VERSION)-darwin-amd64.gz

build-all: \
	bin/grpc-lb-consul-$(VERSION)-linux-amd64 \
	bin/grpc-lb-consul-$(VERSION)-darwin-amd64 \
	bin/grpc-lb-client-$(VERSION)-linux-amd64 \
	bin/grpc-lb-client-$(VERSION)-darwin-amd64

bin/grpc-lb-client-$(VERSION)-linux-amd64: $(SRC)
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/grpc-lb-client
bin/grpc-lb-client-$(VERSION)-darwin-amd64: $(SRC)
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $@ ./cmd/grpc-lb-client

bin/grpc-lb-consul-$(VERSION)-linux-amd64: $(SRC)
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/grpc-lb-consul
bin/grpc-lb-consul-$(VERSION)-darwin-amd64: $(SRC)
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $@ ./cmd/grpc-lb-consul
