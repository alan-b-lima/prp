package user

import (
	"errors"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Repository struct {
	repo map[uuid.UUID]User
}

var ErrUserNotFound = errors.New("user not found")

func NewRepository() Repository {
	repo := Repository{repo: make(map[uuid.UUID]User)}

	uuid1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	uuid2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	repo.repo[uuid1] = User{uuid1, "Alan", "alan"}
	repo.repo[uuid2] = User{uuid2, "Juan", "juan"}

	return repo
}

func (r *Repository) GetAll(req *GetAllRequest) (GetAllResponse, error) {
	if len(r.repo) == 0 {
		return GetAllResponse{}, nil
	}

	res := GetAllResponse(make([]GetResponse, 0, len(r.repo)))
	for _, user := range r.repo {
		res = append(res, _RespondUser(&user))
	}

	return res, nil
}

func (r *Repository) Get(req *GetRequest) (GetResponse, error) {
	user, in := r.repo[req.UUID]
	if !in {
		return GetResponse{}, ErrUserNotFound
	}

	return _RespondUser(&user), nil
}

func _RespondUser(user *User) GetResponse {
	return GetResponse{
		UUID:  user.UUID(),
		Name:  user.Name(),
		Login: user.Login(),
	}
}

func (r *Repository) Create(req *CreateRequest) (CreateResponse, error) {
	user, err := NewUser(&UserScrath{
		Name:  req.Name,
		Login: req.Login,
	})
	if err != nil {
		return CreateResponse{}, err
	}

	r.repo[user.UUID()] = *user
	return CreateResponse{}, nil
}

func (r *Repository) Update(req *UpdateRequest) (UpdateResponse, error) {
	user, in := r.repo[req.UUID]
	if !in {
		return UpdateResponse{}, ErrUserNotFound
	}

	errs := [...]error{
		user.SetName(req.Name),
		user.SetLogin(req.Login),
	}

	if err := errors.Join(errs[:]...); err != nil {
		return UpdateResponse{}, err
	}

	r.repo[req.UUID] = user
	return UpdateResponse{}, nil
}

func (r *Repository) Patch(req *PatchRequest) (PatchResponse, error) {
	user, in := r.repo[req.UUID]
	if !in {
		return PatchResponse{}, ErrUserNotFound
	}

	errs := [2]error{}

	if req.Name != nil  { errs[0] = user.SetName(*req.Name) }
	if req.Login != nil { errs[1] = user.SetLogin(*req.Login) }

	if err := errors.Join(errs[:]...); err != nil {
		return PatchResponse{}, err
	}

	r.repo[req.UUID] = user
	return PatchResponse{}, nil
}

func (r *Repository) Delete(req *DeleteRequest) (DeleteResponse, error) {
	delete(r.repo, req.UUID)
	return DeleteResponse{}, nil
}
