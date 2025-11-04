package session

import (
	"time"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Repository interface {
	Getter
	Creater
}

type Getter interface {
	Get(uuid.UUID) (Entity, error)
}

type Creater interface {
	Create(uuid.UUID, time.Duration) (Entity, error)
}

type Entity struct {
	UUID    uuid.UUID
	User    uuid.UUID
	Expires time.Time
}

func (e Entity) MarshalJSON() ([]byte, error) {
	panic("remember me")
}
