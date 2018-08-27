default: preprocess demo

.PHONY: preprocess
preprocess:
	protoc -I /usr/local/include/ -I. --gogo_out="Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:." options/preprocess.proto  

.PHONY: demo
demo:
	go build
	protoc -I/home/aidyn/go/src -I. -I/usr/local/include --plugin=protoc-gen-preprocess=protoc-gen-preprocess --preprocess_out=./out demo.proto --go_out=./out