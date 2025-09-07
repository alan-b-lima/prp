package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Router struct {
	Repo *Repository
	mux  *http.ServeMux
}

func NewRouter(repo *Repository) Router {
	router := Router{
		Repo: repo,
		mux:  http.NewServeMux(),
	}

	router.mux.HandleFunc(   "GET /users",      router.GetAllHandler)
	router.mux.HandleFunc(   "GET /users/{id}", router.GetHandler)
	router.mux.HandleFunc(  "POST /users",      router.CreateHandler)
	router.mux.HandleFunc(   "PUT /users/{id}", router.UpdateHandler)
	router.mux.HandleFunc( "PATCH /users/{id}", router.PatchHandler)
	router.mux.HandleFunc("DELETE /users/{id}", router.DeleteHandler)

	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := router.Repo.GetAll(&GetAllRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (router *Router) GetHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := GetRequest{UUID: uuid}

	user, err := router.Repo.Get(&req)
	if err != nil {
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

func (router *Router) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = router.Repo.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Router) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	if req.UUID != uuid {
		http.Error(w, "update request body must not contain and uuid field", http.StatusBadRequest)
		return
	}

	_, err = router.Repo.Update(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Router) PatchHandler(w http.ResponseWriter, r *http.Request) {
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

	if req.UUID != uuid {
		http.Error(w, "update request body must not contain and uuid field", http.StatusBadRequest)
		return
	}

	_, err = router.Repo.Patch(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *Router) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := DeleteRequest{UUID: uuid}

	_, err = router.Repo.Delete(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
