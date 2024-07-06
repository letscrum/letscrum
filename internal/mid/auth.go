package mid

import (
	"context"
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func Auth(ctx context.Context) (context.Context, error) {
	enabledAuth := viper.GetBool("server.grpc.interceptors.auth.enabled")
	if !enabledAuth {
		return ctx, nil
	}

	method, _ := grpc.Method(ctx)
	ignoreMethods := viper.GetStringSlice("server.grpc.interceptors.auth.ignoreMethods")

	for _, imethod := range ignoreMethods {
		if strings.Contains(method, imethod) {
			return ctx, nil
		}
	}

	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, fmt.Errorf("auth from md err:%v", err)
	}

	tokenInfo, err := utils.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := context.WithValue(ctx, "token", tokenInfo) // nolint
	return newCtx, nil
}

func AllButHealthZ(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
}
