package user

import (
	"net/http"

	sup "github.com/alan-b-lima/prp/internal/support"
)

type Resource struct {
	http.ServeMux
	Repo Repository
}

func NewResource(repo Repository) *Resource {
	rc := Resource{Repo: repo}

	routes := map[string]sup.RouteFunc{
		"GET /users/":              rc.ListHandler,
		"GET /users/{uuid}":        rc.GetHandler,
		"GET /users/login/{login}": rc.GetByLoginHandler,
		"POST /users/":             rc.CreateHandler,
		"PATCH /users/{uuid}":      rc.PatchHandler,
		"DELETE /users/{uuid}":     rc.DeleteHandler,
	}

	for route, handler := range routes {
		rc.Handle(route, sup.RouteFunc(handler))
	}

	return &rc
}

func (rc *Resource) ListHandler(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	req := ListRequest{Offset: 0, Limit: 10}

	if err := sup.LimitAndOffset(
		query.Get("offset"), query.Get("limit"),
		&req.Offset, &req.Limit,
	); err != nil {
		return err
	}

	res, err := List(rc.Repo, &req)
	if err != nil {
		return err
	}
	if res.Records == nil {
		// avoid "null" encoding, once v2 rolls out,
		// this can be removed
		res.Records = []Response{}
	}

	if err := sup.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) GetHandler(w http.ResponseWriter, r *http.Request) error {
	uuid, err := sup.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := GetRequest{UUID: uuid}
	res, err := Get(rc.Repo, &req)
	if err != nil {
		return err
	}

	if err := sup.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) GetByLoginHandler(w http.ResponseWriter, r *http.Request) error {
	req := GetByLoginRequest{Login: r.PathValue("login")}
	res, err := GetByLogin(rc.Repo, &req)
	if err != nil {
		return err
	}

	if err := sup.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	var req CreateRequest
	if err := sup.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := Create(rc.Repo, &req)
	if err != nil {
		return err
	}

	if err := sup.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) PatchHandler(w http.ResponseWriter, r *http.Request) error {
	uuid, err := sup.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := PatchRequest{UUID: uuid}
	if err := sup.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := Patch(rc.Repo, &req)
	if err != nil {
		return err
	}

	if err := sup.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

func (rc *Resource) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	uuid, err := sup.UUIDFromString(r.PathValue("uuid"))
	if err != nil {
		return err
	}

	req := DeleteRequest{UUID: uuid}
	if err := Delete(rc.Repo, &req); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
