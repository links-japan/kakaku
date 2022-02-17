package handler

import (
	"net/http"

	"github.com/links-japan/kakaku/internal/handler/asset"
	"github.com/links-japan/kakaku/internal/handler/render"

	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New() Server {
	return Server{}
}

type (
	Server struct {
	}
)

func (s Server) HandleRest() http.Handler {
	r := chi.NewRouter()
	r.Use(render.WrapResponse(true))

	r.Get("/assets", asset.HandleAssets())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found 1"))
	})

	return r
}
