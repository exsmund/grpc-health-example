package main

import (
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/exsmund/grpc-health-example/protos/health"
)

var addr = flag.String("addr", "localhost:8001", "The server address")

func main() {
	var grpcServer *grpc.Server
	var mu sync.Mutex

	for {
		go func() {
			time.Sleep(time.Second * 5)
			mu.Lock()
			if grpcServer != nil {
				log.Printf("Stoping gRPC server")
				grpcServer.Stop()
			}
			mu.Unlock()
		}()
		mu.Lock()
		lis, err := net.Listen("tcp", *addr)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("Starting gRPC server %s", lis.Addr().String())
		grpcServer = grpc.NewServer()
		pb.RegisterHealthServer(grpcServer, &healthServer{})
		mu.Unlock()
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
		time.Sleep(time.Second * 5)
	}
}
