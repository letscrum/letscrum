package mid

import (
	"context"
	"fmt"

	"github.com/letscrum/letscrum/pkg/utils"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth  set the token to context.
func Auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, fmt.Errorf("auth from md err:%v", err)
	}

	// todo maybe need to call remote auth service
	tokenInfo, err := utils.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	type ctxKey string
	newCtx := context.WithValue(ctx, ctxKey("tokenInfo"), tokenInfo)

	return newCtx, nil
}
