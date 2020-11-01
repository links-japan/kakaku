package kakaku

import (
	"context"
	kakakupb "github.com/links-japan/kakaku/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

func BTCToJPY() (decimal.Decimal, error) {
	r, err := assetPrice("BTC", "JPY")
	if err != nil {
		return decimal.Zero, nil
	}
	n, _ := decimal.NewFromString(r.Price)
	return n, nil
}

func assetPrice(base, quote string) (*kakakupb.AssetPriceResponse, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("KAKAKU_ADDR", grpc.WithInsecure())
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