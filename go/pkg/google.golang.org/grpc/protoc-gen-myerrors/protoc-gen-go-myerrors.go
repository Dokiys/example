package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const release = "v0.0.1"

var version = flag.Bool("version", false, "print the version and exit")

//go:generate go build -o ./bin/protoc-gen-go-myerrors
func main() {
	flag.Parse()
	if *version {
		fmt.Printf("protoc-gen-go-repo-errors %v\n", release)
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generate(gen, f)
		}
		return nil
	})
}
