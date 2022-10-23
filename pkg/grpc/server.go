package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	port       int
	grpcServer *grpc.Server
}

func New(port int, services []Service) *Server {
	server := grpc.NewServer()
	for _, s := range services {
		s.RegisterGrpcService(server)
	}

	return &Server{
		port:       port,
		grpcServer: server,
	}
}

func (s *Server) Start() <-chan error {
	ch := make(chan error, 1)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		ch <- fmt.Errorf("listen: %w", err)
		return ch
	}

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			ch <- fmt.Errorf("serve: %w", err)
			return
		}
		ch <- nil
	}()

	return ch
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
