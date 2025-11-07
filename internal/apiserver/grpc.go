package apiserver

import (
	"context"
	"net"
	"time"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	pb "github.com/neee333ko/api/proto/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/log"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedCacheServer
	Store    store.Factory
	server   *grpc.Server
	listener net.Listener
}

func (s *GrpcServer) ListSecrets(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	var option metav1.ListOptions
	option.Limit = *req.Limit
	option.Offset = *req.Offset

	sl, err := s.Store.NewSecretStore().List(ctx, &option)
	if err != nil {
		return nil, err
	}

	resp := new(pb.ListSecretsResponse)
	resp.TotalCount = sl.TotalCount
	resp.Items = make([]*pb.SecretInfo, 0)
	for _, s := range sl.Items {
		resp.Items = append(resp.Items, &pb.SecretInfo{
			Name:        s.Name,
			SecretId:    s.SecretID,
			SecretKey:   s.SecretKey,
			Expires:     s.Expires,
			Description: s.Description,
			CreatedAt:   s.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}

func (s *GrpcServer) ListPolicies(ctx context.Context, req *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	var option metav1.ListOptions
	option.Limit = *req.Limit
	option.Offset = *req.Offset

	pl, err := s.Store.NewPolicyStore().List(ctx, &option)
	if err != nil {
		return nil, err
	}

	resp := new(pb.ListPoliciesResponse)
	resp.TotalCount = pl.TotalCount
	resp.Items = make([]*pb.PolicyInfo, 0)
	for _, p := range pl.Items {
		resp.Items = append(resp.Items, &pb.PolicyInfo{
			Name:         p.Name,
			Username:     p.Username,
			PolicyShadow: p.PolicyShadow,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}

func (s *GrpcServer) Run() {
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err.Error())
	}
}

func (s *GrpcServer) Close() {
	log.Info("grpc graceful shutdown...")

	s.server.GracefulStop()
}
