package grpc

import "google.golang.org/grpc"

type Service interface {
	RegisterGrpcService(server grpc.ServiceRegistrar)
}
