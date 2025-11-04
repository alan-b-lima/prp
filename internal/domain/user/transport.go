package user

import (
	"time"

	"github.com/alan-b-lima/prp/pkg/opt"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type (
	ListRequest struct {
		Offset int `json:"-"`
		Limit  int `json:"-"`
	}

	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	GetByLoginRequest struct {
		Login string `json:"-"`
	}

	CreateRequest struct {
		Name     string `json:"name"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	PatchRequest struct {
		UUID     uuid.UUID       `json:"-"`
		Name     opt.Opt[string] `json:"name"`
		Login    opt.Opt[string] `json:"login"`
		Password opt.Opt[string] `json:"password"`
	}

	DeleteRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	AuthRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	ContextRequest struct {
		Session uuid.UUID `json:"-"`
	}
)

type (
	ListResponse struct {
		Offset       int        `json:"offset"`
		Length       int        `json:"length"`
		Records      []Response `json:"records"`
		TotalRecords int        `json:"total_records"`
	}

	AuthResponse struct {
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}

	Response struct {
		UUID  uuid.UUID `json:"uuid"`
		Name  string    `json:"name"`
		Login string    `json:"login"`
		Level string    `json:"level"`
	}
)
