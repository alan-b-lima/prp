package session

import (
	"time"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Getter interface {
	Get(uuid.UUID) (Response, error)
}

type Creater interface {
	Create(uuid.UUID, time.Duration) (Response, error)
}

type Repository interface {
	Getter
	Creater
}
