package user

import (
	"github.com/alan-b-lima/prp/pkg/opt"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type Lister interface {
	List(offset, limit int) (ListEntity, error)
}

type Getter interface {
	Get(uuid uuid.UUID) (Entity, error)
}

type GetterByLogin interface {
	GetByLogin(login string) (Entity, error)
}

type Creater interface {
	Create(name, login, password string) (Entity, error)
}

type Patcher interface {
	Patch(uuid uuid.UUID, name, login, password opt.Opt[string]) (Entity, error)
}

type Deleter interface {
	Delete(uuid uuid.UUID) error
}

type Repository interface {
	Lister
	Getter
	GetterByLogin
	Creater
	Patcher
	Deleter
}

type (
	Entity struct {
		UUID     uuid.UUID
		Name     string
		Login    string
		Password [60]byte
	}

	ListEntity struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}
)
