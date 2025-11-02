package session

import (
	"time"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type (
	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}
	
	CreateRequest struct {
		User   uuid.UUID     `json:"user"`
		MaxAge time.Duration `json:"max_age"`
	}
)

type (
	Response struct {
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}
)
