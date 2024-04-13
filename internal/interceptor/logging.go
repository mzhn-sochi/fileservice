package interceptor

import (
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func LoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		moment := time.Now()
		fmt.Printf("%s - %s\n", moment.String(), info.FullMethod)
		return handler(srv, ss)
	}
}
