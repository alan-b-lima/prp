package user

import (
	"time"

	"github.com/alan-b-lima/prp/internal/domain/session"
	"github.com/alan-b-lima/prp/internal/xerrors"
	"github.com/alan-b-lima/prp/pkg/hash"
)

func List(users Lister, req ListRequest) (ListResponse, error) {
	res, err := users.List(req.Offset, req.Limit)
	if err != nil {
		return ListResponse{}, nil
	}

	ares := ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := 0; i < res.Length; i++ {
		transform(&ares.Records[i], &res.Records[i])
	}

	return ares, nil
}

func Get(users Getter, req GetRequest) (Response, error) {
	res, err := users.Get(req.UUID)
	if err != nil {
		return Response{}, err
	}

	var ares Response
	transform(&ares, &res)
	return ares, nil
}

func GetByLogin(users GetterByLogin, req GetByLoginRequest) (Response, error) {
	res, err := users.GetByLogin(req.Login)
	if err != nil {
		return Response{}, err
	}

	var ares Response
	transform(&ares, &res)
	return ares, nil
}

func Create(users Creater, req CreateRequest) (Response, error) {
	res, err := users.Create(req.Name, req.Login, req.Password)
	if err != nil {
		return Response{}, err
	}

	var ares Response
	transform(&ares, &res)
	return ares, nil
}

func Patch(users Patcher, req PatchRequest) (Response, error) {
	res, err := users.Patch(req.UUID, req.Name, req.Login, req.Password)
	if err != nil {
		return Response{}, err
	}

	var ares Response
	transform(&ares, &res)
	return ares, nil
}

func Delete(users Deleter, req DeleteRequest) error {
	return users.Delete(req.UUID)
}

func Authenticate(users GetterByLogin, sessions session.Creater, req AuthRequest) (AuthResponse, error) {
	res, err := users.GetByLogin(req.Login)
	if err != nil {
		return AuthResponse{}, err
	}

	if !hash.Compare(res.Password[:], []byte(req.Password)) {
		return AuthResponse{}, xerrors.ErrIncorrectPassword
	}

	s, err := session.Create(sessions, session.CreateRequest{
		User:   res.UUID,
		MaxAge: 10 * time.Minute,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	ares := AuthResponse{
		UUID:    s.UUID,
		User:    res.UUID,
		Expires: s.Expires,
	}
	return ares, nil
}

func transform(r *Response, e *Entity) {
	r.UUID = e.UUID
	r.Name = e.Name
	r.Login = e.Login
}
