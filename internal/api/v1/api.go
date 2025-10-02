package api

import (
	"net/http"

	"github.com/alan-b-lima/prp/internal/domain/user"
)

type APIMux struct{ *http.ServeMux }

func NewAPIMux() APIMux {
	mux := APIMux{http.NewServeMux()}

	repo := user.NewRepository()
	router := user.NewResource(repo)

	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.StripPrefix("/api/v1", &router).ServeHTTP(w, r)
	})

	return mux
}
