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

	/*temp*/
	{
		uuid1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
		uuid2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

		repo.repo[uuid1] = User{uuid: uuid1, name: "Alan", login: "alan"}
		repo.repo[uuid2] = User{uuid: uuid2, name: "Vitor", login: "vecto"}
	}

	return repo
}

func (r *Repository) GetAll(req *GetAllRequest) (GetAllResponse, error) {
	if len(r.repo) == 0 {
		return GetAllResponse{}, nil
	}

	res := GetAllResponse(make([]GetResponse, 0, len(r.repo)))
	for _, user := range r.repo {
		res = append(res, respond_user(&user))
	}

	return res, nil
}

func (r *Repository) Get(req *GetRequest) (GetResponse, error) {
	user, in := r.repo[req.UUID]
	if !in {
		return GetResponse{}, ErrUserNotFound
	}

	return respond_user(&user), nil
}

func respond_user(user *User) GetResponse {
	return GetResponse{
		UUID:  user.UUID(),
		Name:  user.Name(),
		Login: user.Login(),
	}
}

func (r *Repository) Create(req *CreateRequest) (CreateResponse, error) {
	user, err := NewUser(&UserScrath{
		Name:     req.Name,
		Login:    req.Login,
		Password: req.Password,
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

	err := errors.Join(
		user.SetName(req.Name),
		user.SetLogin(req.Login),
		user.SetPassword(req.Password),
	)
	if err != nil {
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

	err := errors.Join(
		non_nil_then(req.Name, user.SetName),
		non_nil_then(req.Login, user.SetLogin),
		non_nil_then(req.Password, user.SetPassword),
	)
	if err != nil {
		return PatchResponse{}, err
	}

	r.repo[req.UUID] = user
	return PatchResponse{}, nil
}

func non_nil_then[R any](ptr *R, fn func(R) error) error {
	if ptr != nil {
		return fn(*ptr)
	}

	return nil
}

func (r *Repository) Delete(req *DeleteRequest) (DeleteResponse, error) {
	delete(r.repo, req.UUID)
	return DeleteResponse{}, nil
}
