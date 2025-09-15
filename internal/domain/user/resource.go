package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Resource struct {
	Repo *Repository
	mux  *http.ServeMux
}

func NewResource(repo *Repository) Resource {
	router := Resource{
		Repo: repo,
		mux:  http.NewServeMux(),
	}

	router.mux.HandleFunc("GET /users", router.GetAllHandler)
	router.mux.HandleFunc("GET /users/{id}", router.GetHandler)
	router.mux.HandleFunc("POST /users", router.CreateHandler)
	router.mux.HandleFunc("PUT /users/{id}", router.UpdateHandler)
	router.mux.HandleFunc("PATCH /users/{id}", router.PatchHandler)
	router.mux.HandleFunc("DELETE /users/{id}", router.DeleteHandler)

	return router
}

func (router *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Resource) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	req := GetAllRequest{Limit: 100}
	var res GetAllResponse

	if err := router.Repo.GetAll(&res, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (router *Resource) GetHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := GetRequest{UUID: uuid}
	var user GetResponse

	if err := router.Repo.Get(&user, &req); err != nil {
		if err == ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (router *Resource) CreateHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req CreateRequest
	err = json.Unmarshal(buf, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := router.Repo.Create(nil, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Resource) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := UpdateRequest{UUID: uuid}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := router.Repo.Update(nil, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Resource) PatchHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := PatchRequest{UUID: uuid}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = router.Repo.Patch(nil, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Resource) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := DeleteRequest{UUID: uuid}

	if err := router.Repo.Delete(nil, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
