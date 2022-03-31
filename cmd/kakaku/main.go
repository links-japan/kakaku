package main

import (
	"os"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/links-japan/kakaku/internal/client"
	"github.com/links-japan/kakaku/internal/config"
	"github.com/links-japan/kakaku/internal/kakaku"
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

	go startWorker()
	select {}
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
		log.Panic("start worker", err)
	}

	for _, asset := range lst {
		oracle := kakaku.NewOracle(clients, assets, &cfg.Oracle)
		go Run(oracle, assets, asset.Base, asset.Quote)
	}
}

func Run(oracle *kakaku.Oracle, assets *store.AssetStore, base, quote string) {
	for {
		if err := kakaku.UpdateAssetPrice(oracle, assets, base, quote); err != nil {
			log.Errorln("update asset price error", err)
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
