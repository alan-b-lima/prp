package userrepo

import (
	"cmp"
	"sync"

	"github.com/alan-b-lima/prp/internal/domain/user"
	"github.com/alan-b-lima/prp/pkg/errors"
	"github.com/alan-b-lima/prp/pkg/opt"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Map struct {
	uuidIndex  map[uuid.UUID]int
	loginIndex map[string]int

	repo []user.User
	mu   sync.RWMutex
}

func NewMap() user.Repository {
	repo := Map{
		uuidIndex:  make(map[uuid.UUID]int),
		loginIndex: make(map[string]int),
	}

	{
		repo.Create(&user.CreateRequest{Name: "Alan Lima", Login: "alan-b-lima", Password: "12345678"})
		repo.Create(&user.CreateRequest{Name: "Juan Ferreira", Login: "juanzinho_bs", Password: "12345678"})
		repo.Create(&user.CreateRequest{Name: "Luan Filipe", Login: "lf-carvalho", Password: "12345678"})
		repo.Create(&user.CreateRequest{Name: "Mateus Oliveira", Login: "mateuzinhodelasoficial2013", Password: "12345678"})
		repo.Create(&user.CreateRequest{Name: "Vitor Mozer", Login: "vecto", Password: "12345678"})
	}

	return &repo
}

func (m *Map) List(req *user.ListRequest) (user.ListResponse, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, req.Offset, len(m.repo))
	hi := clamp(0, req.Offset+req.Limit, len(m.repo))

	if lo >= hi {
		return user.ListResponse{TotalRecords: len(m.repo)}, nil
	}

	res := make([]user.Response, hi-lo)
	for i := range m.repo[lo:hi] {
		transform(&res[i], &m.repo[i])
	}

	return user.ListResponse{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) Get(req *user.GetRequest) (user.Response, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[req.UUID]
	if !in {
		return user.Response{}, user.ErrUserNotFound
	}

	var res user.Response
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) GetByLogin(req *user.GetByLoginRequest) (user.Response, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.loginIndex[req.Login]
	if !in {
		return user.Response{}, user.ErrUserNotFound
	}

	var res user.Response
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) Create(req *user.CreateRequest) (user.Response, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	u, err := user.NewUser(&user.Scratch{
		Name:     req.Name,
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return user.Response{}, err
	}

	if _, in := m.loginIndex[u.Login()]; in {
		return user.Response{}, user.ErrLoginTaken
	}

	m.uuidIndex[u.UUID()] = len(m.repo)
	m.loginIndex[u.Login()] = len(m.repo)
	m.repo = append(m.repo, u)

	var res user.Response
	transform(&res, &u)
	return res, nil
}

func (m *Map) Patch(req *user.PatchRequest) (user.Response, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[req.UUID]
	if !in {
		return user.Response{}, user.ErrUserNotFound
	}

	u := m.repo[index]
	oldLogin := u.Login()

	err := errors.Join(
		some_then(req.Name, u.SetName),
		some_then(req.Login, u.SetLogin),
		some_then(req.Password, u.SetPassword),
	)
	if err != nil {
		return user.Response{}, err
	}

	m.repo[index] = u

	delete(m.loginIndex, oldLogin)
	m.loginIndex[u.Login()] = index

	var res user.Response
	transform(&res, &u)
	return res, nil
}

func (m *Map) Delete(req *user.DeleteRequest) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[req.UUID]
	if !in {
		return nil
	}

	u := &m.repo[index]

	delete(m.uuidIndex, u.UUID())
	delete(m.loginIndex, u.Login())

	if len(m.repo) == 1 {
		m.repo = m.repo[:0]
		return nil
	}

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]

	return nil
}

func some_then[F any](val opt.Opt[F], fn func(F) error) error {
	if val.Some {
		return fn(val.Val)
	}
	return nil
}

func transform(r *user.Response, u *user.User) {
	r.UUID = u.UUID()
	r.Name = u.Name()
	r.Login = u.Login()
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
