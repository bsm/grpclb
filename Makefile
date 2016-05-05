default: vet errcheck test

test:
	go test ./...

vet:
	go tool vet -printf=false -composites=false $(wildcard *.go)

errcheck:
	errcheck -ignoretests -ignore 'Close' $$(go list ./...)

proto: $(patsubst %.proto,%.pb.go,$(wildcard */*.proto))

.PHONY: default test vet errcheck all proto

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<
