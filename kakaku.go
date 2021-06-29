package kakaku

import (
	"context"
	kakakupb "github.com/links-japan/kakaku/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"os"
	"time"
)

func PriceWithTime(base, quote string) (decimal.Decimal, time.Time, error) {
	if _, ok := os.LookupEnv("KAKAKU_FAKE_DATA"); !ok {
		return decimal.NewFromInt(1000000), time.Now(), nil
	}
	r, err := AssetPrice(base, quote)
	if err != nil {
		return decimal.Zero, time.Time{}, err
	}
	n, _ := decimal.NewFromString(r.Price)
	return n, r.Timestamp.AsTime(), nil
}

func AssetPrice(base, quote string) (*kakakupb.AssetPriceResponse, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(os.Getenv("KAKAKU_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := kakakupb.NewCheckinServiceClient(conn)
	req, err := c.AssetPrice(
		context.Background(),
		&kakakupb.AssetPriceRequest{
			Base:  base,
			Quote: quote,
		},
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}
