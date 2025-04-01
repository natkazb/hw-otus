package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"   //nolint
	"google.golang.org/grpc/status" //nolint
)

const UNKNOWN = "?"

func logging(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)

		ip := UNKNOWN
		peerInfo, ok := peer.FromContext(ctx)
		if ok {
			ip = peerInfo.Addr.String()
		}

		userAgent := UNKNOWN
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			userAgent = md.Get("user-agent")[0]
		}

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		log.Info(fmt.Sprintf("%s  %s  %d  %s  %v",
			ip,
			info.FullMethod,
			statusCode,
			userAgent,
			time.Since(start)))

		return resp, err
	}
}
