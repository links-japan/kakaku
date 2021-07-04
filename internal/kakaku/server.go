package kakaku

import (
	"context"
	"fmt"
	"github.com/links-japan/kakaku/internal/store"
	kakakupb "github.com/links-japan/kakaku/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	assets *store.AssetStore
}

func NewServer(assets *store.AssetStore) *Server {
	return &Server{assets: assets}
}

func (s *Server) AssetPrice(ctx context.Context, req *kakakupb.AssetPriceRequest) (*kakakupb.AssetPriceResponse, error) {
	base, quote := req.Base, req.Quote

	all, err := s.assets.ListAll()
	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	assets := filterAssets(all)
	price, ok := search(assets, base, quote)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not found")
	}

	return &kakakupb.AssetPriceResponse{
		Base:      req.Base,
		Quote:     req.Quote,
		Price:     price.String(),
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

func filterAssets(all []*store.Asset) []*store.Asset {
	var assets []*store.Asset
	for _, a := range all {
		if a.Term > 0 {
			assets = append(assets, a)
		}
	}
	return assets
}

type Edge struct {
	From string
	To   string
	Val  decimal.Decimal
}

func buildGraph(assets []*store.Asset) *map[string][]Edge {
	graph := make(map[string][]Edge)
	for _, a := range assets {
		graph[a.Base] = append(graph[a.Base], Edge{
			From: a.Base,
			To:   a.Quote,
			Val:  a.Price,
		})
		graph[a.Quote] = append(graph[a.Quote], Edge{
			From: a.Quote,
			To:   a.Base,
			Val:  decimal.NewFromInt(1).Div(a.Price),
		})
	}
	return &graph
}

func dfs(graph *map[string][]Edge, visit *map[string]int, path *[]Edge, edge Edge, target string) bool {
	if (*visit)[edge.To] > 0 {
		return false
	}
	(*visit)[edge.To] = 1

	*path = append(*path, edge)

	if edge.To == target {
		return true
	}

	for _, e := range (*graph)[edge.To] {
		if dfs(graph, visit, path, e, target) {
			return true
		}
	}

	*path = (*path)[:len(*path)-1]
	return false
}

func search(assets []*store.Asset, base, quote string) (decimal.Decimal, bool) {
	graph := buildGraph(assets)
	visit := make(map[string]int)
	var path []Edge

	start := Edge{
		From: "",
		To:   base,
		Val:  decimal.NewFromInt(1),
	}

	dfs(graph, &visit, &path, start, quote)
	fmt.Println(path)
	if len(path) < 2 {
		return decimal.Zero, false
	}

	res := decimal.NewFromInt(1)
	for _, e := range path {
		res = res.Mul(e.Val)
	}
	return res, true
}
