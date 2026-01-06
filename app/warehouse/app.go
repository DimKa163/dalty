package warehouse

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimKa163/dalty/pkg/graph"
	"github.com/DimKa163/dalty/pkg/proto"

	"github.com/DimKa163/dalty/internal/logging"
	"github.com/DimKa163/dalty/internal/warehouse/core"
	"github.com/DimKa163/dalty/internal/warehouse/persistence"
	"github.com/DimKa163/dalty/internal/warehouse/server"
	"github.com/DimKa163/dalty/internal/warehouse/server/interceptor"
	"github.com/DimKa163/dalty/internal/warehouse/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type ServiceContainer struct {
	PathService         *usecase.PathService
	PgPool              *pgxpool.Pool
	GrpcServer          *grpc.Server
	GrpcPathServer      proto.Binder
	WarehouseRepository core.WarehouseRepository
	GraphContext        *graph.GraphContext
	binders             []proto.Binder
}

func (s *ServiceContainer) GetBinders() []proto.Binder {
	return s.binders
}

type Server struct {
	Config *Config
	*ServiceContainer
	proto.ServerImpl
}

func NewServer(config *Config) *Server {
	container := &ServiceContainer{
		binders: make([]proto.Binder, 0),
	}
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
	s.ServerImpl = proto.NewGRPCServer[*ServiceContainer](listener, addGrpcServer(), s.ServiceContainer)
	s.GraphContext = addGraphContext()
	s.PgPool, err = addPgPool(s.Config.Database)
	if err != nil {
		return err
	}
	s.WarehouseRepository = addWarehouseRepository(s.PgPool)
	s.PathService = addPathService(s.WarehouseRepository, s.GraphContext)
	s.binders = append(s.binders, addGrpcPathServer(s.PathService))
	return nil
}

func (s *Server) AddLogging() error {
	return logging.InitializeLogging(&logging.LogConfiguration{
		Builders: map[string]logging.CoreBuilder{
			"file":    logging.NewFileBuilder("D:\\logs\\warehouse.log", zap.NewProductionEncoderConfig(), zapcore.InfoLevel),
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

func addGraphContext() *graph.GraphContext {
	return graph.NewGraphContext()
}

func addPathService(repository core.WarehouseRepository,
	graphContext *graph.GraphContext) *usecase.PathService {
	return usecase.NewPathService(repository, core.NewPathFinder(graphContext), graphContext)
}

func addGrpcPathServer(appService *usecase.PathService) *server.PathServer {
	return server.NewPathServer(appService)
}
