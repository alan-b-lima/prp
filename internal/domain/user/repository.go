package user

import (
	"errors"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Repository interface {
	GetAll(*GetAllResponse, *GetAllRequest) error
	Get(*GetResponse, *GetRequest) error
	Create(*CreateResponse, *CreateRequest) error
	Update(*UpdateResponse, *UpdateRequest) error
	Patch(*PatchResponse, *PatchRequest) error
	Delete(*DeleteResponse, *DeleteRequest) error
}

type MapRepository struct {
	repo map[uuid.UUID]User
}

var (
	ErrUserNotFound       = errors.New("user(repo): not found")
	ErrLoginAlreadyExists = errors.New("user(repo): login already in use")
)

func NewRepository() Repository {
	repo := MapRepository{repo: make(map[uuid.UUID]User)}

	/*temp*/
	{
		uuid1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
		uuid2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

		repo.repo[uuid1] = User{uuid: uuid1, name: "Alan", login: "alan"}
		repo.repo[uuid2] = User{uuid: uuid2, name: "Vitor", login: "vecto"}
	}

	return &repo
}

func (r *MapRepository) GetAll(res *GetAllResponse, req *GetAllRequest) error {
	if len(r.repo) == 0 {
		*res = nil
		return nil
	}

	*res = make(GetAllResponse, len(r.repo))

	var i int
	for _, user := range r.repo {
		transform(&(*res)[i], &user)
		i++
	}

	return nil
}

func (r *MapRepository) Get(res *GetResponse, req *GetRequest) error {
	user, in := r.repo[req.UUID]
	if !in {
		return ErrUserNotFound
	}

	transform(res, &user)
	return nil
}

func (r *MapRepository) Create(res *CreateResponse, req *CreateRequest) error {
	user, err := NewUser(&UserScrath{
		Name:     req.Name,
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return insert(r, user)
}

func (r *MapRepository) Update(res *UpdateResponse, req *UpdateRequest) error {
	user, in := r.repo[req.UUID]
	if !in {
		return ErrUserNotFound
	}

	err := errors.Join(
		user.SetName(req.Name),
		user.SetLogin(req.Login),
		user.SetPassword(req.Password),
	)
	if err != nil {
		return err
	}

	return insert(r, &user)
}

func (r *MapRepository) Patch(res *PatchResponse, req *PatchRequest) error {
	user, in := r.repo[req.UUID]
	if !in {
		return ErrUserNotFound
	}

	err := errors.Join(
		non_nil_then(req.Name, user.SetName),
		non_nil_then(req.Login, user.SetLogin),
		non_nil_then(req.Password, user.SetPassword),
	)
	if err != nil {
		return err
	}

	return insert(r, &user)
}

func (r *MapRepository) Delete(res *DeleteResponse, req *DeleteRequest) error {
	delete(r.repo, req.UUID)
	return nil
}

func transform(res *GetResponse, user *User) {
	res.UUID = user.UUID()
	res.Name = user.Name()
	res.Login = user.Login()
}

func insert(r *MapRepository, user *User) error {
	for _, u := range r.repo {
		if user.Login() == u.Login() {
			if user.UUID() == u.UUID() {
				continue
			}

			return ErrLoginAlreadyExists
		}
	}

	r.repo[user.UUID()] = *user
	return nil
}

func non_nil_then[R any](ptr *R, fn func(R) error) error {
	if ptr != nil {
		return fn(*ptr)
	}

	return nil
}
