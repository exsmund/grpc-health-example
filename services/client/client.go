package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/exsmund/grpc-health-example/protos/health"
)

var serverAddr = flag.String("serveraddr", "localhost:8001", "The address to connect to")

func healthCheck(ctx context.Context, conn *grpc.ClientConn) {
	c := pb.NewHealthClient(conn)
	for {
		select {
		case <-ctx.Done():
			// The context is over, stop processing results
			log.Print("Context is over")
			return
		case <-time.After(time.Second * time.Duration(rand.Intn(4))):
			// Repeat the healthcheck after random timeout
			res, err := c.Check(ctx, &pb.HealthCheckRequest{})
			if err != nil {
				log.Print("Could not reach the server")
			} else {
				log.Printf("Status: %s", res.GetStatus())
			}
		}
	}
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var wg sync.WaitGroup

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctxChild, cancel := context.WithCancel(ctx)
		defer cancel()
		healthCheck(ctxChild, conn)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctxChild, cancel := context.WithCancel(ctx)
		defer cancel()
		healthCheck(ctxChild, conn)
	}()

	wg.Wait()
}
