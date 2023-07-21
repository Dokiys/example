.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 &&  \
	go install github.com/gogo/protobuf/protoc-gen-gofast@v1.3.2 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 &&  \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.11.3 &&  \
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors@v2.5.1 && \
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.8

HELLO_PROTO_FILES=$(shell find ./helloproto/* -name *.proto)
HELLO_PROTO_GRPC_FILES=$(shell find ./hellogrpc/* -name *.proto)
PROTO_GATEWAY_FILES=$(shell find ./gateway/* -name *.proto)

.PHONY: proto
proto:
	protoc --proto_path=. \
		--descriptor_set_out=helloproto/descriptor.desc --go_out=../ \
		$(PROTO_FILES)

.PHONY: gofast_proto
gofast_proto:
	protoc --proto_path=. \
		--gofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:../ \
		$(PROTO_FILES)

# GRPC
.PHONY: proto_grpc
proto_grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(HELLO_PROTO_GRPC_FILES)

# gRPC Gateway
.PHONY: proto_gateway
OPENAPI_PATH=./gateway
proto_gateway:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./ \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        --openapi_out=${OPENAPI_PATH} \
        $(HELLO_PROTO_GRPC_FILES)

# protoc-gen-myerrors
include ./protoc-gen-myerrors/Makefile