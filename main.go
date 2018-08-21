package main

import (
	"github.com/atorgayev/protoc-gen-preprocess/preprocessor"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	p := preprocessor.NewPreprocessor()
	p.Overwrite()
	resp := command.GeneratePlugin(req, p, ".preprocess.pb.go")
	command.Write(resp)
}
