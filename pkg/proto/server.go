package proto

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"reflect"
)

type Servicer interface {
	GetBinders() []Binder
}

type ServerImpl interface {
	ListenAndServe() error
	Map()
	Shutdown(ctx context.Context) error
}

type GRPCServer[T Servicer] struct {
	servicer Servicer
	listener net.Listener
	*grpc.Server
}

func NewGRPCServer[T Servicer](listener net.Listener, server *grpc.Server, servicer T) *GRPCServer[T] {
	return &GRPCServer[T]{
		Server:   server,
		listener: listener,
		servicer: servicer,
	}
}

func (gs *GRPCServer[T]) ListenAndServe() error {
	fmt.Printf("starting gRPC Server on %s ...\n", gs.listener.Addr().String())
	return gs.Serve(gs.listener)
}

func (gs *GRPCServer[T]) Map() {
	binders := gs.servicer.GetBinders()
	for _, b := range binders {
		el := reflect.TypeOf(b).Elem()
		fmt.Printf("mapping %s\n", el.Name())
		b.Bind(gs.Server)
	}
}

func (gs *GRPCServer[T]) Shutdown(ctx context.Context) error {
	gs.GracefulStop()
	fmt.Println("server shutdown gracefully")
	return nil
}
