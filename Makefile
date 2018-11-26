VERSION=v$(strip $(shell cat .version))

default: vet test

test:
	go test ./...

vet:
	go vet ./...

.PHONY: test vet

# ---------------------------------------------------------------------

proto: proto.go proto.python proto.ruby

touch-proto:
	touch $(wildcard */*.proto)

force-proto: touch-proto proto

proto.go: $(patsubst %.proto,%.pb.go,$(wildcard */*.proto))
%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<

proto.python: $(patsubst %.proto,%_pb2.py,$(wildcard */*.proto))
%_pb2.py: %.proto
	@if ! python -c 'import grpc_tools'; then echo 'Run "pip install grpcio-tools" required to generate python code'; exit 1; fi

	python -m grpc_tools.protoc -I. --python_out=python --grpc_python_out=python $<

proto.ruby: $(patsubst %.proto,%_pb.rb,$(wildcard */*.proto))
%_pb.rb: %.proto
	bundle exec grpc_tools_ruby_protoc --ruby_out=ruby/lib --grpc_out=ruby/lib $<

.PHONY: proto touch-proto force-proto

# ---------------------------------------------------------------------

TARGET_PKG=$(patsubst cmd/%/main.go,bin/%,$(wildcard cmd/grpc-lb-*/main.go))
TARGET_OS=linux darwin
TARGET_ARCH=amd64 386
TARGETS=$(foreach pkg,$(TARGET_PKG),$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(pkg)-$(os)-$(arch))))

all: $(TARGETS)
dist: $(foreach t,$(TARGETS),$(t).zip)

bin/grpc-lb-%.zip: bin/grpc-lb-%
	zip -j $@ $<
bin/grpc-lb-%: $(shell find . -name '*.go')
	@mkdir -p $(dir $@)
	$(eval os := $(word 4, $(subst -, ,$@)))
	$(eval arch := $(word 5, $(subst -, ,$@)))
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -ldflags "-X main.version=$(VERSION)" -o $@ $(patsubst bin/%-$(os)-$(arch),cmd/%/main.go,$@)

.PHONY: all dist
