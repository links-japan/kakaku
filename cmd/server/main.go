package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/drone/signal"
	"github.com/go-chi/chi"
	"github.com/links-japan/kakaku/internal/config"
	"github.com/links-japan/kakaku/internal/handler"
	"github.com/links-japan/kakaku/internal/handler/hc"
	"github.com/links-japan/kakaku/internal/store"
	"github.com/links-japan/log"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

var (
	cfg config.Config
)

func main() {
	initConfig()
	log.Init()

	if err := store.Connect(cfg.DB.Dsn); err != nil {
		log.Fatal(err)
	}
	store.Conn().AutoMigrate(&store.Asset{})
	go startServer()
	select {}
}

func startServer() {
	ctx := context.Background()
	mux := chi.NewMux()
	// hc
	{
		mux.Mount("/", hc.HandleHc())
	}
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
		log.Debug("shutdown server...")

		// create context with timeout
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			log.WithError(err).Error("graceful shutdown server failed")
		}

		close(done)
	})

	log.Infoln("serve at", addr)
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("server aborted")
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("KAKAKU_CONFIG_PATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	delta, err := decimal.NewFromString(cfg.Oracle.PriceDeltaStr)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Oracle.PriceDelta = delta
}
