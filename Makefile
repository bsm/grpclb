SRC:=$(shell find . -name '*.go' -not -path '*vendor*')
PKG:=$(shell glide nv)
PROTO:=$(patsubst %.proto,%.pb.go,$(wildcard */*.proto))

default: vet test

test:
	go test $(PKG)

vet:
	go tool vet -printf=false -composites=false $(SRC)

proto: $(PROTO)

touch-proto:
	touch $(wildcard */*.proto)

force-proto: touch-proto proto

.PHONY: default test vet all proto

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<
