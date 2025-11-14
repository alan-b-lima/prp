package api

import (
	"net/http"

	sessionrepo "github.com/alan-b-lima/prp/internal/domain/session/repository"
	userrepo "github.com/alan-b-lima/prp/internal/domain/user/repository"
	users "github.com/alan-b-lima/prp/internal/domain/user/resource"
)

type router struct{ http.ServeMux }

func New() http.Handler {
	var r router

	var (
		sessionsRepo = sessionrepo.NewMap()
		usersRepo    = userrepo.NewMap()
	)

	users := users.New(usersRepo, sessionsRepo)

	r.Handle("/api/v1/users/", http.StripPrefix("/api/v1", users))
	return &r
}
