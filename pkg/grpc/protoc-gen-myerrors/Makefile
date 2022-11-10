ERROR_PROTO_FILES=$(shell find ./proto/* -name *.proto)
.PHONY: genproto
genproto:
	go install . && protoc -I=./proto \
	  -I=../third_party \
	  --go_out=paths=source_relative:./gen \
	  --go-myerrors_out=paths=source_relative:./gen \
	  ${ERROR_PROTO_FILES}
