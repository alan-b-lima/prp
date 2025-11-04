package user

import (
	"github.com/alan-b-lima/prp/internal/auth"
	"github.com/alan-b-lima/prp/internal/domain/session"
	"github.com/alan-b-lima/prp/internal/xerrors"
)

type Service struct {
	Repo     Repository
	Sessions session.Repository
}

func NewService(users Repository, sessions session.Repository) *Service {
	return &Service{
		Repo:     users,
		Sessions: sessions,
	}
}

var (
	PermAdmin      = auth.Permission(auth.Admin)
	PermGeneral    = auth.Permission(auth.Admin, auth.User)
	PermPermissive = auth.Permission(auth.Admin, auth.User, auth.Unlogged)
)

func (s *Service) List(ctx auth.Context, req ListRequest) (ListResponse, error) {
	if l, c := ctx.Level(), PermAdmin; !c.Authorize(l) {
		return ListResponse{}, xerrors.ErrUnauthorizedUser.New(l, c)
	}

	return List(s.Repo, req)
}

func (s *Service) Get(ctx auth.Context, req GetRequest) (Response, error) {
	if ctx.User() == req.UUID {
		goto Do
	}

	if l, c := ctx.Level(), PermAdmin; !c.Authorize(l) {
		return Response{}, xerrors.ErrUnauthorizedUser.New(l, c)
	}

Do:
	return Get(s.Repo, req)
}

func (s *Service) GetByLogin(ctx auth.Context, req GetByLoginRequest) (Response, error) {
	res, err := GetByLogin(s.Repo, req)
	if err != nil && res.Login == req.Login {
		goto Do
	}

	if l, c := ctx.Level(), PermAdmin; !c.Authorize(l) {
		return Response{}, xerrors.ErrUnauthorizedUser.New(l, c)
	}

Do:
	return res, err
}

func (s *Service) Create(ctx auth.Context, req CreateRequest) (Response, error) {
	if l, c := ctx.Level(), PermPermissive; !c.Authorize(l) {
		return Response{}, xerrors.ErrUnauthorizedUser.New(l, c)
	}

	return Create(s.Repo, req)
}

func (s *Service) Patch(ctx auth.Context, req PatchRequest) (Response, error) {
	if ctx.User() == req.UUID {
		goto Do
	}

	if l, c := ctx.Level(), PermAdmin; !c.Authorize(l) {
		return Response{}, xerrors.ErrUnauthorizedUser.New(l, c)
	}

Do:
	return Patch(s.Repo, req)
}

func (s *Service) Delete(ctx auth.Context, req DeleteRequest) error {
	if ctx.User() == req.UUID {
		goto Do
	}

	if l, c := ctx.Level(), PermAdmin; !c.Authorize(l) {
		return xerrors.ErrUnauthorizedUser.New(l, c)
	}

Do:
	return Delete(s.Repo, req)
}

func (s *Service) Authenticate(req AuthRequest) (AuthResponse, error) {
	return Authenticate(s.Repo, s.Sessions, req)
}

func (s *Service) Context(req ContextRequest) (auth.Context, error) {
	return Context(s.Repo, s.Sessions, req)
}
