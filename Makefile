SRC=$(shell find . -name '*.go' -not -path '*vendor*')
PKG=$(shell go list ./... | grep -v 'vendor')
VERSION=v0.3.1

GOPROTO=$(patsubst %.proto,%.pb.go,$(wildcard */*.proto))
PYPROTO=$(patsubst %.proto,%_pb2.py,$(wildcard */*.proto))
RBPROTO=$(patsubst %.proto,%_pb.rb,$(wildcard */*.proto))

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


proto: proto.go proto.python proto.ruby

touch-proto:
	touch $(wildcard */*.proto)

force-proto: touch-proto proto

build: $(TARGETS)
dist: $(ARCHIVES)

.PHONY: default test vet all proto force-proto build dist

proto.go: $(GOPROTO)
%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<

proto.python: $(PYPROTO)
%_pb2.py: %.proto
	@if ! python -c 'import grpc_tools'; then echo 'Run "pip install grpcio-tools" required to generate python code'; exit 1; fi

	python -m grpc_tools.protoc -I. --python_out=python --grpc_python_out=python $<

proto.ruby: $(RBPROTO)
%_pb.rb: %.proto
	protoc --ruby_out=ruby/lib --grpc_out=ruby/lib --plugin=protoc-gen-grpc=`which grpc_ruby_plugin` $<

bin/grpc-lb-%.zip: bin/grpc-lb-%
	zip -j $@ $<
bin/grpc-lb-%: $(SRC)
	@mkdir -p $(dir $@)
	$(eval os := $(word 4, $(subst -, ,$@)))
	$(eval arch := $(word 5, $(subst -, ,$@)))
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -o $@ $(patsubst bin/%-$(os)-$(arch),cmd/%/main.go,$@)
