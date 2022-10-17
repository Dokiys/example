.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 \
	&& go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors@latest \
	&& go install .


ERROR_PROTO_FILES=$(shell find ./proto/* -name *.proto)
.PHONY: genproto
genproto:
	# --go-errors_out=paths=source_relative:./gen
	go install . && protoc -I=./proto \
	  -I=./third_party/kratos/errors \
	  --go_out=paths=source_relative:./gen \
	  --go-myerrors_out=paths=source_relative:./gen \
	  ${ERROR_PROTO_FILES}
