package use

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
)

func Auth2() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        fmt.Println(info.FullMethod)
        //md, ok := metadata.FromIncomingContext(ctx)
        //if !ok {
        //    return errors.New("无Token认证信息")
        //}
        var newCtx context.Context
        fmt.Println("auth")
        return handler(newCtx, req)
    }
}
//
//func Use(ctx context.Context, info *grpc.UnaryServerInfo) error {
//    fmt.Println(info.FullMethod)
//    md, ok := metadata.FromIncomingContext(ctx)
//    if !ok {
//        return errors.New("无Token认证信息")
//    }
//    var (
//        appid  string
//        appkey string
//    )
//    if val, ok := md["appid"]; ok {
//        appid = val[0]
//    }
//    if val, ok := md["appkey"]; ok {
//        appkey = val[0]
//    }
//    if appid != "101010" || appkey != "i am key" {
//        return errors.New("Token认证信息无效")
//    }
//    return nil
//}

//
//func Auth() grpc.StreamServerInterceptor {
//    return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//        wrapped := grpc_middleware.WrapServerStream(stream)
//        fmt.Println("auth")
//        return handler(srv, wrapped)
//    }
//}
//
//func Auth2() grpc.UnaryServerInterceptor {
//    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//        var newCtx context.Context
//        fmt.Println("auth")
//        return handler(newCtx, req)
//    }
//}
