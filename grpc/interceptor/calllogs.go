package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// Interceptor Interceptor
type Interceptor struct {
}

// UnaryServerInterceptor 拦截器
func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("新调用. 调用信息: %+v 请求信息: %+v", info, req)
	resp, err := handler(ctx, req)
	log.Printf("调用结束. 结果: %+v", resp)
	return resp, err
}

// StreamServerInterceptor 拦截器
func (i *Interceptor) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("新流调用. 调用信息: %+v ", info)
	err := handler(srv, ss)
	log.Printf("流调用结束. 结果: %v", err)
	return err
}

// NewCalllogs 创建拦截器
func NewCalllogs() []grpc.ServerOption {
	interceptor := new(Interceptor)
	unary := grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor)
	stream := grpc.StreamInterceptor(interceptor.StreamServerInterceptor)
	result := []grpc.ServerOption{unary, stream}
	return result
}
