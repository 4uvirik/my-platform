package observability

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start).Seconds()

		method := info.FullMethod

		GRPCRequests.WithLabelValues(method).Inc()
		GRPCLatency.WithLabelValues(method).Observe(elapsed)

		if err != nil {
			st, _ := status.FromError(err)
			_ = st
			GRPCErrors.WithLabelValues(method).Inc()
		}

		return resp, err
	}
}
