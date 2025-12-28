package graph

import (
	"context"
	"github.com/DimKa163/dalty/internal/graph/core"
	"github.com/DimKa163/dalty/internal/graph/persistence"
	"github.com/DimKa163/dalty/internal/graph/server"
	"github.com/DimKa163/dalty/internal/graph/server/interceptor"
	"github.com/DimKa163/dalty/internal/graph/usecase"
	"github.com/DimKa163/dalty/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type ServerImpl interface {
	ListenAndServe() error
	Map()
	Shutdown(ctx context.Context) error
}
type ServiceContainer struct {
	PathService         *usecase.PathService
	PgPool              *pgxpool.Pool
	GrpcServer          *grpc.Server
	GrpcPathServer      *server.PathServer
	WarehouseRepository core.WarehouseRepository
	GraphContext        *core.GraphContext
}
type Server struct {
	Config *Config
	*ServiceContainer
	ServerImpl
}

func NewServer(config *Config) *Server {
	container := &ServiceContainer{}
	return &Server{
		Config:           config,
		ServiceContainer: container,
	}
}

func (s *Server) AddServices() error {
	var err error
	listener, err := net.Listen("tcp", s.Config.Addr)
	if err != nil {
		return err
	}
	s.ServerImpl = NewGRPCServer(listener, addGrpcServer(), s.ServiceContainer)
	s.GraphContext = addGraphContext()
	s.PgPool, err = addPgPool(s.Config.Database)
	if err != nil {
		return err
	}
	s.WarehouseRepository = addWarehouseRepository(s.PgPool)
	s.PathService = addPathService(s.WarehouseRepository, s.GraphContext)
	s.GrpcPathServer = addGrpcPathServer(s.PathService)
	return nil
}

func (s *Server) AddLogging() error {
	return logging.InitializeLogging(&logging.LogConfiguration{
		Builders: map[string]logging.CoreBuilder{
			"file":    logging.NewFileBuilder("D:\\logs\\graph.log", zap.NewProductionEncoderConfig(), zapcore.InfoLevel),
			"console": logging.NewConsoleBuilder(zap.NewDevelopmentEncoderConfig(), zapcore.DebugLevel),
		},
	})
}

func (s *Server) Map() {
	s.ServerImpl.Map()
}

func (s *Server) Run() error {
	logger := logging.GetLogger().Sugar()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	if err := s.PathService.UpdateGraph(ctx); err != nil {
		logger.Errorf("PathService.UpdateGraph err: %v", err)
		return err
	}
	s.addSyscallObserver(ctx)
	return s.ListenAndServe()
}

func (s *Server) addSyscallObserver(ctx context.Context) {
	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		logger := logging.Logger(timeoutCtx)
		logger.Info("graceful shutdown")
		_ = s.Shutdown(timeoutCtx)
	}()
}
func addGrpcServer() *grpc.Server {
	return grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryServerLoggingInterceptor()))
}

func addPgPool(database string) (*pgxpool.Pool, error) {
	pg, err := pgxpool.New(context.Background(), database)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func addWarehouseRepository(pool *pgxpool.Pool) core.WarehouseRepository {
	return persistence.NewWarehouseRepository(pool)
}

func addGraphContext() *core.GraphContext {
	return core.NewGraphContext()
}

func addPathService(repository core.WarehouseRepository,
	graphContext *core.GraphContext) *usecase.PathService {
	return usecase.NewPathService(repository, graphContext)
}

func addGrpcPathServer(appService *usecase.PathService) *server.PathServer {
	return server.NewPathServer(appService)
}
