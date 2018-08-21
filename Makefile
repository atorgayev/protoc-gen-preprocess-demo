default:
	go build
	protoc --plugin=protoc-gen-preprocess=protoc-gen-preprocess --preprocess_out=./out demo.proto