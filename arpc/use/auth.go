package use

import (
    "context"
    "fmt"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    "google.golang.org/grpc"
)

func Auth() grpc.StreamServerInterceptor {
    return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        wrapped := grpc_middleware.WrapServerStream(stream)
        fmt.Println("auth")
        return handler(srv, wrapped)
    }
}

func Auth2() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        var newCtx context.Context
        fmt.Println("auth")
        return handler(newCtx, req)
    }
}