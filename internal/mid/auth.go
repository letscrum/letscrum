package mid

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Auth(ctx context.Context) (context.Context, error) {
	return ctx, nil
	//token, err := auth.AuthFromMD(ctx, "bearer")
	//if err != nil {
	//	return nil, fmt.Errorf("auth from md err:%v", err)
	//}
	//
	//tokenInfo, err := utils.ParseToken(token)
	//if err != nil {
	//	return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	//}
	//
	//type ctxKey string
	//newCtx := context.WithValue(ctx, ctxKey("tokenInfo"), tokenInfo)
	//
	//return newCtx, nil
}

func AllButHealthZ(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
}
