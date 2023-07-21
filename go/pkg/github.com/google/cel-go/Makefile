.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30

PROTO_FILES=$(shell find ./* -name *.proto)

.PHONY: proto
proto:
	protoc --proto_path=./example  --go_out=../ $(PROTO_FILES)

