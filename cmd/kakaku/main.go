package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/links-japan/kakaku/internal/client"
	"github.com/links-japan/kakaku/internal/config"
	"github.com/links-japan/kakaku/internal/kakaku"
	"github.com/links-japan/kakaku/internal/store"
	kakakupb "github.com/links-japan/kakaku/pb"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"time"
)

var (
	cfg config.Config
)

func main() {
	initConfig()
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := store.Connect(cfg.DB.Dsn); err != nil {
		log.Fatal(err)
	}
	store.Conn().AutoMigrate(&store.Asset{})

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
	// TODO: set worker tasks in config
	base := "BTC"
	quote := "JPY"

	keystore := mixin.Keystore{
		ClientID:   cfg.Mixin.ClientID,
		SessionID:  cfg.Mixin.SessionID,
		PrivateKey: cfg.Mixin.PrivateKey,
		PinToken:   cfg.Mixin.PinToken,
	}

	clients := []client.Client{
		client.NewMixinClient(&keystore),
		client.NewLiquidClient(),
		client.NewCoinBaseClient(),
	}
	assets := store.NewAssetStore()

	asset := store.Asset{
		Base:   base,
		Quote:  quote,
		Price:  decimal.Zero,
		Source: "",
		Term:   0,
	}
	_ = assets.FirstOrCreate(&asset)

	oracle := kakaku.NewOracle(clients, assets, &cfg.Oracle)

	for {
		if err := kakaku.UpdateAssetPrice(oracle, assets, base, quote); err != nil {
			logrus.Errorln("update asset price error", err)
		}
		time.Sleep(cfg.Worker.TermTimeout)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("KAKAKU_CONFIG_PATH"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatal(err)
	}

	delta, err := decimal.NewFromString(cfg.Oracle.PriceDeltaStr)
	if err != nil {
		logrus.Fatal(err)
	}
	cfg.Oracle.PriceDelta = delta
}
