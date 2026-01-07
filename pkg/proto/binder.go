package proto

import "google.golang.org/grpc"

type Binder interface {
	Bind(server *grpc.Server)
}
