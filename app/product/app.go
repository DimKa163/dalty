package product

import (
	"context"
	"github.com/DimKa163/dalty/internal/product/server/interceptor"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimKa163/dalty/internal/logging"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/internal/product/persistence"
	"github.com/DimKa163/dalty/internal/product/server"
	"github.com/DimKa163/dalty/internal/product/usecase"
	"github.com/DimKa163/dalty/pkg/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type ServiceContainer struct {
	PgPool             *pgxpool.Pool
	ProductRepository  core.ProductRepository
	RelationRepository core.RelationRepository
	GrpcServer         *grpc.Server
	ProductService     *usecase.ProductService
	binders            []proto.Binder
	ProductServer      proto.Binder
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
	s.PgPool, err = addPgPool(s.Config.Database)
	if err != nil {
		return err
	}
	s.ProductRepository = addProductRepository(s.PgPool)
	s.ProductService = addProductService(s.ProductRepository)
	s.RelationRepository = addRelationRepository(s.PgPool)
	s.binders = append(s.binders, server.NewProductServer(s.ProductService),
		server.NewSpecificationServer(usecase.NewSpecificationService(s.ProductRepository, s.RelationRepository)))
	s.ServerImpl = proto.NewGRPCServer[*ServiceContainer](listener, addGrpcServer(), s.ServiceContainer)
	return nil
}

func (s *Server) AddLogging() error {
	return logging.InitializeLogging(&logging.LogConfiguration{
		Builders: map[string]logging.CoreBuilder{
			"file":    logging.NewFileBuilder("D:\\logs\\product.log", zap.NewProductionEncoderConfig(), zapcore.InfoLevel),
			"console": logging.NewConsoleBuilder(zap.NewDevelopmentEncoderConfig(), zapcore.DebugLevel),
		},
	})
}

func (s *Server) Map() {
	s.ServerImpl.Map()
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
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
func addProductRepository(pool *pgxpool.Pool) core.ProductRepository {
	return persistence.NewProductRepository(pool)
}

func addRelationRepository(pool *pgxpool.Pool) core.RelationRepository {
	return persistence.NewRelationRepository(pool)
}

func addProductService(productService core.ProductRepository) *usecase.ProductService {
	return usecase.NewProductService(productService)
}
