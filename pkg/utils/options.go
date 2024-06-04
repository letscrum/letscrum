package utils

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

type Options struct {
	Network    string
	GRPCAddr   string
	HTTPAddr   string
	OpenAPIDir string
	Mux        []runtime.ServeMuxOption
}
