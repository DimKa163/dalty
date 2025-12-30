package interceptor

import (
	"context"
	"time"

	"github.com/DimKa163/dalty/internal/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := logging.Logger(ctx)
		mData, _ := metadata.FromIncomingContext(ctx)
		logger = logger.With(zap.Any("server", info.Server))
		logger = logger.With(zap.Any("method", info.FullMethod))
		if logger.Level() == zap.DebugLevel {
			logger = logger.With(zap.Any("body", req))
		}
		for k, v := range mData {
			logger = logger.With(zap.Any(k, v))
		}
		logger.Info("got incoming request")
		ctx = logging.SetLogger(ctx, logger)
		startTime := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(startTime)
		if err != nil {
			logger.Error("request failed", zap.Error(err))
		}
		logger = logger.With(zap.Duration("elapsed", elapsed)).With(zap.Any("response", resp))
		logger.Info("request processed")
		return resp, err
	}
}
