package kakaku

import (
	"context"
	"github.com/links-japan/kakaku/internal/store"
	kakakupb "github.com/links-japan/kakaku/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct{}

func (s *Server) AssetPrice(ctx context.Context, req *kakakupb.AssetPriceRequest) (*kakakupb.AssetPriceResponse, error) {
	base, quote := req.Base, req.Quote
	if base != "BTC" || quote != "JPY" {
		return nil, status.Error(codes.Unimplemented, "unimplemented")
	}

	var asset store.Asset
	err := store.Conn().Where("base = ? AND quote = ?", base, quote).First(&asset).Error
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &kakakupb.AssetPriceResponse{
		Base:      req.Base,
		Quote:     req.Quote,
		Price:     asset.Price.String(),
		Timestamp: timestamppb.New(asset.UpdatedAt),
	}, nil
}
