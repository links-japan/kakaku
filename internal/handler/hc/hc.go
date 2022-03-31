package hc

import (
	"net/http"

	"github.com/links-japan/kakaku/internal/handler/render"
)

func HandleHc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, map[string]interface{}{"status": "OK"})
		return
	}
}
