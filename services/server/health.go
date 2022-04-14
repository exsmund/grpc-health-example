package main

import (
	"context"

	pb "github.com/exsmund/grpc-health-example/protos/health"
)

type healthServer struct {
	pb.UnimplementedHealthServer
}

func (s *healthServer) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(req *pb.HealthCheckRequest, stream pb.Health_WatchServer) error {
	stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})
	return nil
}
