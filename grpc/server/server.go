package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/a2ush/sample-grpc-server-with-appmesh-xray/rpc"
	"github.com/go-redis/redis"
)

type TimeManageServer struct{}

var counter int

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen; %v", err)
	}

	// Add logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(zapLogger),
			),
		),
	)

	rpc.RegisterTimeManageServer(
		server,
		&TimeManageServer{},
	)

	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}

func (s *TimeManageServer) ConvertTime(
	ctx context.Context,
	req *rpc.ClientRequest,
) (*rpc.ServerResponse, error) {

	timezone, _ := time.LoadLocation("UTC")
	switch req.TimezoneFormat {
	case rpc.Timezone_UTC:
		timezone, _ = time.LoadLocation("UTC")
	case rpc.Timezone_PST:
		timezone, _ = time.LoadLocation("America/New_York")
	case rpc.Timezone_JST:
		timezone, _ = time.LoadLocation("Asia/Tokyo")
	}

	convert_time := time.Now().In(timezone)
	log.Printf("req.TimezoneFormat: %v, return : %v", req.TimezoneFormat, convert_time.Format(time.RFC3339))

	// record the result to redis
	client := redis.NewClient(&redis.Options{
		Addr: "redis-server.redis.svc.cluster.local:6379",
	})
	counter++
	err := client.Set(string(counter), convert_time, 0).Err()
	if err != nil {
		panic(err)
	}

	return &rpc.ServerResponse{
		ConvertTime: string(convert_time.Format(time.RFC3339)),
	}, nil
}
