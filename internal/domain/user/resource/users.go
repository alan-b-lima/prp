package users

import (
	"net/http"

	"github.com/alan-b-lima/prp/internal/domain/session"
	"github.com/alan-b-lima/prp/internal/domain/user"
	"github.com/alan-b-lima/prp/internal/support"
	"github.com/alan-b-lima/prp/internal/xerrors"
	uuidpkg "github.com/alan-b-lima/prp/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Users    user.Repository
	Sessions session.Repository
}

func New(users user.Repository, sessions session.Repository) *Resource {
	rc := Resource{
		Users:    users,
		Sessions: sessions,
	}

	routes := map[string]support.RouteFunc{
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
		rc.Handle(route, support.RouteFunc(handler))
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	req := user.ListRequest{Offset: 0, Limit: 10}

	if err := support.LimitAndOffset(
		query.Get("offset"), query.Get("limit"),
		&req.Offset, &req.Limit,
	); err != nil {
		return err
	}

	res, err := user.List(rc.Users, req)
	if err != nil {
		return err
	}
	if res.Records == nil {
		// avoid "null" encoding, once v2 rolls out,
		// this can be removed
		res.Records = []user.Response{}
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) error {
	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := user.GetRequest{UUID: uuid}
	res, err := user.Get(rc.Users, req)
	if err != nil {
		return err
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) GetByLogin(w http.ResponseWriter, r *http.Request) error {
	req := user.GetByLoginRequest{Login: r.PathValue("login")}
	res, err := user.GetByLogin(rc.Users, req)
	if err != nil {
		return err
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) error {
	var req user.CreateRequest
	if err := support.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := user.Create(rc.Users, req)
	if err != nil {
		return err
	}

	if err := support.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) error {
	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := user.PatchRequest{UUID: uuid}
	if err := support.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := user.Patch(rc.Users, req)
	if err != nil {
		return err
	}

	if err := support.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) error {
	uuid, err := support.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := user.DeleteRequest{UUID: uuid}
	if err := user.Delete(rc.Users, req); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (rc *Resource) Authenticate(w http.ResponseWriter, r *http.Request) error {
	var req user.AuthRequest
	if err := support.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := user.Authenticate(rc.Users, rc.Sessions, req)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    res.UUID.String(),
		Expires:  res.Expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // should change for HTTPS in the future, I think
	})

	if err := support.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) error {
	s, err := r.Cookie("session")
	if err != nil {
		return xerrors.ErrUnauthenticatedUser
	}

	uuid, err := support.UUIDFromString(s.Value)
	if err != nil {
		return err
	}

	res, err := session.Get(rc.Sessions, session.GetRequest{
		UUID: uuid,
	})
	if err != nil {
		return err
	}

	ures, err := user.Get(rc.Users, user.GetRequest{
		UUID: res.User,
	})
	if err != nil {
		return err
	}

	ares := struct {
		Session uuidpkg.UUID `json:"session"`
		user.Response
	}{
		Session:  res.UUID,
		Response: ures,
	}

	if err := support.EncodeJSON(&ares, http.StatusCreated, w, r); err != nil {
		return err
	}

	return nil
}
