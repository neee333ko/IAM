package apiserver

import (
	"net"

	"github.com/neee333ko/log"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Address string
	*grpc.Server
}

func (s *GRPCServer) Run() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	go func() {
		log.Info("grpc service starts...")

		if err := s.Server.Serve(listener); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return nil
}

func (s *GRPCServer) Close() {
	s.Server.GracefulStop()

	log.Info("grpc service stops.")
}
