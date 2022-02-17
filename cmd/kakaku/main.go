package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/drone/signal"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/go-chi/chi"
	"github.com/links-japan/kakaku/internal/client"
	"github.com/links-japan/kakaku/internal/config"
	"github.com/links-japan/kakaku/internal/handler"
	"github.com/links-japan/kakaku/internal/kakaku"
	"github.com/links-japan/kakaku/internal/store"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	ctx := context.Background()
	mux := chi.NewMux()
	// rpc & api v1 & ws
	{
		svr := handler.New()

		// api v1
		restHandler := svr.HandleRest()
		mux.Mount("/api", restHandler)
	}

	// launch server
	addr := fmt.Sprintf(":%d", 8080)

	svr := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	done := make(chan struct{}, 1)
	ctx = signal.WithContextFunc(ctx, func() {
		logrus.Debug("shutdown server...")

		// create context with timeout
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			logrus.WithError(err).Error("graceful shutdown server failed")
		}

		close(done)
	})

	logrus.Infoln("serve at", addr)
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("server aborted")
	}
}

func startWorker() {

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

	lst, err := assets.ListVariable()
	if err != nil {
		logrus.Panic("start worker", err)
	}

	for _, asset := range lst {
		oracle := kakaku.NewOracle(clients, assets, &cfg.Oracle)
		go Run(oracle, assets, asset.Base, asset.Quote)
	}
}

func Run(oracle *kakaku.Oracle, assets *store.AssetStore, base, quote string) {
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
