package graph

import (
	"context"
	"net"

	"github.com/DimKa163/dalty/internal/logging"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	services *ServiceContainer
	listener net.Listener
	*grpc.Server
}

func NewGRPCServer(listener net.Listener, server *grpc.Server, services *ServiceContainer) *GRPCServer {
	return &GRPCServer{
		Server:   server,
		listener: listener,
		services: services,
	}
}
func (gs *GRPCServer) ListenAndServe() error {
	logger := logging.GetLogger()
	loggerSugar := logger.Sugar()
	loggerSugar.Infof("Listening on %s", gs.listener.Addr())
	return gs.Serve(gs.listener)
}

func (gs *GRPCServer) Map() {
	gs.services.GrpcPathServer.Bind(gs.Server)
}

func (gs *GRPCServer) Shutdown(ctx context.Context) error {
	logger := logging.Logger(ctx)
	gs.GracefulStop()
	logger.Info("server shutdown gracefully")
	return nil
}
