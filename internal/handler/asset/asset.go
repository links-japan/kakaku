package asset

import (
	"net/http"

	"github.com/links-japan/kakaku/internal/handler/render"
	"github.com/links-japan/kakaku/internal/store"
	"github.com/twitchtv/twirp"
)

func HandleAssets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assets := store.NewAssetStore()
		all, err := assets.ListAll()
		if err != nil {
			render.Error(w, twirp.NewError(twirp.Internal, "handle assets err"))
			return
		}
		render.JSON(w, all)
		return
	}
}
