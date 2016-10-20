SRC=$(shell find . -name '*.go' -not -path '*vendor*')
PKG=$(shell glide nv)
PROTO=$(patsubst %.proto,%.pb.go,$(wildcard */*.proto))
VERSION=v0.3.1

TARGET_PKG=$(patsubst cmd/%/main.go,bin/%,$(wildcard cmd/grpc-lb-*/main.go))
TARGET_OS=linux darwin
TARGET_ARCH=amd64 386
TARGETS=$(foreach pkg,$(TARGET_PKG),$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(pkg)-$(os)-$(arch))))
ARCHIVES=$(foreach t,$(TARGETS),$(t).zip)

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)

proto: $(PROTO)

touch-proto:
	touch $(wildcard */*.proto)

force-proto: touch-proto proto

build: $(TARGETS)
dist: $(ARCHIVES)

.PHONY: default test vet all proto force-proto build dist

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<
bin/grpc-lb-%.zip: bin/grpc-lb-%
	zip -j $@ $<
bin/grpc-lb-%: $(SRC)
	@mkdir -p $(dir $@)
	$(eval os := $(word 4, $(subst -, ,$@)))
	$(eval arch := $(word 5, $(subst -, ,$@)))
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -o $@ $(patsubst bin/%-$(os)-$(arch),cmd/%/main.go,$@)
