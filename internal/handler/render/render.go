package render

import (
	"encoding/json"
	"net/http"

	"github.com/twitchtv/twirp"
)

type H map[string]interface{}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func Text(w http.ResponseWriter, t string) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(t))
}

func Error(w http.ResponseWriter, err twirp.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]interface{}{
		"code": err.Code(),
		"msg":  err.Msg(),
	})
}

func Html(w http.ResponseWriter, t string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(t))
}
