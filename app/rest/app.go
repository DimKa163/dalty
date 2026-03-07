package rest

import (
	"context"
	"github.com/DimKa163/dalty/internal/logging"
	"github.com/DimKa163/dalty/internal/rest/core"
	"github.com/DimKa163/dalty/internal/rest/persistence"
	"github.com/DimKa163/dalty/internal/rest/server"
	"github.com/DimKa163/dalty/internal/rest/usecase"
	"github.com/DimKa163/dalty/pkg/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	PgPool         *pgxpool.Pool
	GrpcServer     *grpc.Server
	binders        []proto.Binder
	RestRepository core.RestRepository
	RestService    *usecase.RestService
	RestsServer    proto.Binder
}

func (s *ServiceContainer) GetBinders() []proto.Binder {
	return s.binders
}

type Server struct {
	Context  context.Context
	Config   *Config
	listener net.Listener
	*ServiceContainer
	proto.ServerImpl
}

func NewServer(config *Config) (*Server, error) {
	ctx := context.Background()

	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}
	pg, err := pgxpool.New(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}
	container := &ServiceContainer{
		binders: make([]proto.Binder, 0),
		PgPool:  pg,
	}
	return &Server{
		Context:          ctx,
		Config:           config,
		listener:         listener,
		ServiceContainer: container,
	}, nil
}

func (s *Server) AddServices() {
	s.RestRepository = persistence.NewRestRepository(s.PgPool)
	s.RestService = usecase.NewRestService(s.RestRepository)
	s.binders = append(s.binders, server.NewRestServer(s.RestService))
	s.GrpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(proto.UnaryServerLoggingInterceptor()))
	s.ServerImpl = proto.NewGRPCServer[*ServiceContainer](s.listener, s.GrpcServer, s.ServiceContainer)
}

func (s *Server) AddLogging() error {
	return logging.InitializeLogging(&logging.LogConfiguration{
		Builders: map[string]logging.CoreBuilder{
			"file":    logging.NewFileBuilder("D:\\logs\\rest.log", zap.NewProductionEncoderConfig(), zapcore.InfoLevel),
			"console": logging.NewConsoleBuilder(zap.NewDevelopmentEncoderConfig(), zapcore.DebugLevel),
		},
	})
}

func (s *Server) Map() {
	s.ServerImpl.Map()
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(s.Context, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	log := logging.GetLogger()
	log = log.With(zap.String("addr", s.Config.Addr))
	logging.SetLogger(ctx, log)
	reflection.Register(s.GrpcServer)
	s.addSyscallObserver(ctx)
	logSug := log.Sugar()
	if s.Config.UseProfiling {
		s.pprof(ctx)
	}
	logSug.Infof("with args: %v", s.Config)
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

func (s *Server) pprof(ctx context.Context) {
	go func() {
		log := logging.Logger(ctx)
		log.Info("pprof start")
		if err := http.ListenAndServe(s.Config.PProfAddr, nil); err != nil {
			panic(err)
		}
		<-ctx.Done()
	}()
}
