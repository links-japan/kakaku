package kakaku

import (
	"context"
	"github.com/links-japan/kakaku/internal/store"
	kakakupb "github.com/links-japan/kakaku/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Server struct {
	assets *store.AssetStore
}

func NewServer(assets *store.AssetStore) *Server {
	return &Server{assets: assets}
}

func (s *Server) AssetPrice(ctx context.Context, req *kakakupb.AssetPriceRequest) (*kakakupb.AssetPriceResponse, error) {
	base, quote := req.Base, req.Quote

	var asset store.Asset
	err := s.assets.Find(&asset, base, quote)
	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Term 0's data is invalid
	if asset.Term < 1 {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &kakakupb.AssetPriceResponse{
		Base:      req.Base,
		Quote:     req.Quote,
		Price:     asset.Price.String(),
		Timestamp: timestamppb.New(asset.UpdatedAt),
	}, nil
}
