#!/bin/bash

build() {
	files=${1}/*.proto
	# echo ${files}
    protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		${files}
}

clean() {
	files=${1}/*.pb.go
	rm ${files}
}


clean ${1}
build ${1}
