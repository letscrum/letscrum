.PHONY: build clean tool lint help
all: build

# $Env:GOOS = "e" $Env:GOOS = "darwin"
# export GOOS=linux
build:
	go build -o dist/letscrum ./cmd/letscrum/

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

run:
	./dist/letscrum server

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"

.PHONY: api_gen api_dep_install api_clean
api_dep_install:
	go env -w GOPROXY=https://goproxy.cn,direct
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/jstemmer/go-junit-report@latest
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest
	go install github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts@latest
	go install github.com/rakyll/statik@latest

api_gen:
	protoc -I . -I third_party \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--grpc-gateway_out=paths=source_relative:. \
		--openapiv2_out=logtostderr=true:. \
		--grpc-gateway-ts_out=paths=source_relative:./dist/sdk/ \
		api/general/v1/common.proto \
		api/general/v1/letscrum.proto \
		api/letscrum/v1/letscrum.proto \
		api/project/v1/project.proto \
		api/project/v1/sprint.proto \
		api/item/v1/epic.proto \
		api/item/v1/feature.proto \
		api/item/v1/work_item.proto \
		api/item/v1/task.proto \
		api/user/v1/user.proto
	cp api/letscrum/v1/letscrum.swagger.json docs/swagger-ui/letscrum.swagger.json

api_clean:
	rm -f api/*/*/*.pb.go api/*/*/*.pb.gw.go api/*/*/*.swagger.json api/*/*/*.pb.validate.go
	rm -rf dist/sdk/*
	rm -rf docs/swagger-ui/*.swagger.json
