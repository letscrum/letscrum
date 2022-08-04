.PHONY: build clean tool lint help
all: build

build:
	# $Env:GOOS = "linux"
	go build -o letscrum ./cmd

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"

.PHONY: api_gen api_dep_install api_clean
api_dep_install:
	go env
	go env -w GOPROXY=https://goproxy.cn,direct
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/jstemmer/go-junit-report@latest
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

api_gen:
	protoc -I . -I third_party \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--grpc-gateway_out=paths=source_relative:. \
		--openapiv2_out=logtostderr=true:. \
		--grpc-gateway-ts_out=paths=source_relative:./dist/sdk/ \
		--validate_out=lang=go,paths=source_relative:. \
		apis/general/v1/common.proto apis/general/v1/letscrum.proto apis/letscrum/v1/letscrum.proto apis/project/v1/project.proto apis/user/v1/user.proto apis/project/v1/sprint.proto \

api_clean:
	rm -f apis/*/*/*.pb.go apis/*/*/*.pb.gw.go apis/*/*/*.swagger.json apis/*/*/*.pb.validate.go
	rm -rf dist/sdk/*
	rm -rf third_party/swagger-ui/*.swagger.json
