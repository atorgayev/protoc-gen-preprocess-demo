default: preprocess demo

.PHONY: preprocess
preprocess:
	protoc --gogo_out=. ./options/preprocess.proto  -I /usr/local/include/ -I.

.PHONY: demo
demo:
	go build
	protoc -I/home/aidyn/go/src -I. -I/usr/local/include --plugin=protoc-gen-preprocess=protoc-gen-preprocess --preprocess_out=./out demo.proto