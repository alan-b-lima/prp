package api

import (
	"net/http"

	"github.com/alan-b-lima/prp/internal/domain/user"
)

func NewAPIMux() http.Handler {
	mux := http.NewServeMux()

	repo := user.NewRepository()
	router := user.NewResource(&repo)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", &router))

	return mux
}
