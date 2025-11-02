package sessionrepo

import (
	"sync"
	"time"

	"github.com/alan-b-lima/prp/internal/domain/session"
	"github.com/alan-b-lima/prp/internal/xerrors"
	"github.com/alan-b-lima/prp/pkg/heap"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Map struct {
	uuidIndex   map[uuid.UUID]int
	userIndex   map[uuid.UUID]int
	expiresHeap sleepqueue

	repo []session.Session
	mu   sync.RWMutex
}

func NewMap() session.Repository {
	repo := Map{
		uuidIndex: make(map[uuid.UUID]int),
		userIndex: make(map[uuid.UUID]int),
		expiresHeap: sleepqueue{
			new:    make(chan es, 32),
			cancel: make(chan struct{}, 1),
		},
	}

	go repo.flush()

	return &repo
}

func (m *Map) Get(uuid uuid.UUID) (session.Response, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return session.Response{}, xerrors.ErrSessionNotFound
	}

	s := m.repo[index]
	if time.Now().After(s.Expires()) {
		return session.Response{}, xerrors.ErrSessionNotFound
	}

	var res session.Response
	transform(&res, &m.repo[index])
	return res, nil
}

func (m *Map) Create(user uuid.UUID, maxAge time.Duration) (session.Response, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	if index, in := m.userIndex[user]; in {
		s := m.repo[index]
		m.delete(s.UUID())
	}

	s, err := session.New(&session.Scratch{
		User:   user,
		MaxAge: 1 * time.Minute,
	})
	if err != nil {
		return session.Response{}, err
	}

	m.uuidIndex[s.UUID()] = len(m.repo)
	m.userIndex[s.User()] = len(m.repo)
	m.repo = append(m.repo, s)

	m.expiresHeap.new <- es{s.UUID(), s.Expires()}

	var res session.Response
	transform(&res, &s)
	return res, nil
}

func (m *Map) delete(uuid uuid.UUID) error {
	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	s := &m.repo[index]

	delete(m.uuidIndex, s.UUID())
	delete(m.userIndex, s.User())

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]
	return nil
}

func transform(r *session.Response, s *session.Session) {
	r.UUID = s.UUID()
	r.User = s.User()
	r.Expires = s.Expires()
}

func (m *Map) flush() {
	h := m.expiresHeap

	for {
		var after <-chan time.Time
		if h.heap.Len() > 0 {
			delay := time.Until(h.heap.Peek().expires)
			after = time.After(delay)
		}

		select {
		case <-h.cancel:
			return

		case es := <-h.new:
			h.heap.Push(es)

		case <-after:
			es := h.heap.Pop()

			m.mu.Lock()
			m.delete(es.session)
			m.mu.Unlock()
		}
	}
}

type sleepqueue struct {
	heap   heap.Heap[es]
	new    chan es
	cancel chan struct{}
}

type es struct {
	session uuid.UUID
	expires time.Time
}

func (o0 es) Less(o1 es) bool { return o0.expires.Before(o1.expires) }
