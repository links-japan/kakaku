package main

import (
	"github.com/links-japan/kakaku/internal/kakaku"
	kakakupb "github.com/links-japan/kakaku/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"time"
)

func main() {
	if err := kakaku.Connect(); err != nil {
		log.Fatal(err)
	}
	kakaku.Conn().AutoMigrate(&kakaku.Asset{})

	go startWorker()
	go startServer()

	select {}
}

func startServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s, grpcServer := kakaku.Server{}, grpc.NewServer()
	kakakupb.RegisterCheckinServiceServer(grpcServer, &s)

	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func startWorker() {
	for {
		if err := kakaku.UpdateAssetPrice(); err != nil {
			logrus.Errorln("update asset price error", err)
		}
		time.Sleep(time.Minute)
	}
}
