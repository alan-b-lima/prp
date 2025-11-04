package userrepo

import (
	"cmp"
	"sync"

	"github.com/alan-b-lima/prp/internal/auth"
	"github.com/alan-b-lima/prp/internal/domain/user"
	"github.com/alan-b-lima/prp/internal/xerrors"
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
		repo.Create("Alan Lima", "alan-b-lima", "12345678", auth.Admin)
		repo.Create("Juan Ferreira", "juanzinho_bs", "12345678", auth.User)
		repo.Create("Luan Filipe", "lf-carvalho", "12345678", auth.User)
		repo.Create("Mateus Oliveira", "mateuzinhodelasoficial2013", "12345678", auth.User)
		repo.Create("Vitor Mozer", "vecto", "12345678", auth.User)
	}

	return &repo
}

func (m *Map) List(offset, limit int) (user.ListEntity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, offset, len(m.repo))
	hi := clamp(0, offset+limit, len(m.repo))

	if lo >= hi {
		return user.ListEntity{TotalRecords: len(m.repo)}, nil
	}

	res := make([]user.Entity, hi-lo)
	for i := range m.repo[lo:hi] {
		transform(&res[i], &m.repo[i])
	}

	return user.ListEntity{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) Get(uuid uuid.UUID) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	var res user.Entity
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) GetByLogin(login string) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.loginIndex[login]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	var res user.Entity
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) Create(name, login, password string, level auth.Level) (user.Entity, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	u, err := user.New(&user.Scratch{
		Name:     name,
		Login:    login,
		Password: password,
		Level:    level,
	})
	if err != nil {
		return user.Entity{}, err
	}

	if _, in := m.loginIndex[u.Login()]; in {
		return user.Entity{}, xerrors.ErrLoginTaken
	}

	m.uuidIndex[u.UUID()] = len(m.repo)
	m.loginIndex[u.Login()] = len(m.repo)
	m.repo = append(m.repo, u)

	var res user.Entity
	transform(&res, &u)
	return res, nil
}

func (m *Map) Patch(uuid uuid.UUID, name, login, password opt.Opt[string]) (user.Entity, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return user.Entity{}, xerrors.ErrUserNotFound
	}

	u := m.repo[index]
	oldLogin := u.Login()

	err := errors.Join(
		some_then(name, u.SetName),
		some_then(login, u.SetLogin),
		some_then(password, u.SetPassword),
	)
	if err != nil {
		return user.Entity{}, err
	}

	m.repo[index] = u

	delete(m.loginIndex, oldLogin)
	m.loginIndex[u.Login()] = index

	var res user.Entity
	transform(&res, &u)
	return res, nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	u := &m.repo[index]

	delete(m.uuidIndex, u.UUID())
	delete(m.loginIndex, u.Login())

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

func transform(r *user.Entity, u *user.User) {
	r.UUID = u.UUID()
	r.Name = u.Name()
	r.Login = u.Login()
	r.Password = u.Password()
	r.Level = u.Level()
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
