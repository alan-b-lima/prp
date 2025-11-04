package users

import (
	"net/http"

	"github.com/alan-b-lima/prp/internal/auth"
	"github.com/alan-b-lima/prp/internal/domain/session"
	"github.com/alan-b-lima/prp/internal/domain/user"
	"github.com/alan-b-lima/prp/internal/support"
	"github.com/alan-b-lima/prp/internal/xerrors"
	"github.com/alan-b-lima/prp/pkg/errors"
)

const _SessionCookie = "session"

type Resource struct {
	http.ServeMux
	Users user.Service
}

func New(users user.Repository, sessions session.Repository) *Resource {
	rc := Resource{
		Users: *user.NewService(users, sessions),
	}

	routes := map[string]http.HandlerFunc{
		"GET /users/":              rc.List,
		"GET /users/{uuid}":        rc.Get,
		"GET /users/login/{login}": rc.GetByLogin,
		"POST /users/":             rc.Create,
		"PATCH /users/{uuid}":      rc.Patch,
		"DELETE /users/{uuid}":     rc.Delete,
		"POST /users/auth/":        rc.Authenticate,
		"GET /users/me/":           rc.Me,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	query := r.URL.Query()
	req := user.ListRequest{Offset: 0, Limit: 10}

	if err := support.LimitAndOffset(
		query.Get("offset"), query.Get("limit"),
		&req.Offset, &req.Limit,
	); err != nil {
		support.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.List(ctx, req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}
	if res.Records == nil {
		// avoid "null" encoding, once v2 rolls out,
		// this can be removed
		res.Records = []user.Response{}
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	req := user.GetRequest{UUID: uuid}
	res, err := rc.Users.Get(ctx, req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) GetByLogin(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	req := user.GetByLoginRequest{Login: r.PathValue("login")}
	res, err := rc.Users.GetByLogin(ctx, req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	var req user.CreateRequest
	if err := support.DecodeJSON(&req, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Create(ctx, req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	if err := support.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	req := user.PatchRequest{UUID: uuid}
	if err := support.DecodeJSON(&req, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Patch(ctx, req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, err := rc.session(w, r)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	req := user.DeleteRequest{UUID: uuid}
	if err := rc.Users.Delete(ctx, req); err != nil {
		support.WriteJsonError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rc *Resource) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req user.AuthRequest
	if err := support.DecodeJSON(&req, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Authenticate(req)
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     _SessionCookie,
		Value:    res.UUID.String(),
		Expires:  res.Expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // should change for HTTPS in the future, I think
	})

	if err := support.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) {
	session, err := support.SessionCookie(_SessionCookie, w, r)
	if err != nil {
		support.WriteJsonError(w, xerrors.ErrUnauthenticatedUser)
		return
	}
	
	ctx, err := rc.Users.Context(user.ContextRequest{Session: session})
	if err != nil {
		support.WriteJsonError(w, xerrors.ErrUnauthenticatedUser)
		return
	}

	ures, err := rc.Users.Get(ctx, user.GetRequest{UUID: ctx.User()})
	if err != nil {
		support.WriteJsonError(w, err)
		return
	}

	if err := support.EncodeJSON(&ures, http.StatusOK, w, r); err != nil {
		support.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) session(w http.ResponseWriter, r *http.Request) (auth.Context, error) {
	session, err := support.SessionCookie(_SessionCookie, w, r)
	if err != nil {
		return auth.NewUnlogged(), nil
	}

	ctx, err := rc.Users.Context(user.ContextRequest{Session: session})
	if err, ok := errors.AsType[*errors.Error](err); ok && err.Kind.IsClient() {
		return auth.NewUnlogged(), nil
	}
	if err != nil {
		return auth.NewUnlogged(), err
	}

	return ctx, err
}
