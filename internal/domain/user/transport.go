package user

import "github.com/alan-b-lima/prp/pkg/uuid"

type NoContent struct{}

type (
	GetAllRequest struct {
		Limit int `json:"limit"`
	}

	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	CreateRequest struct {
		Name     string `json:"name"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	UpdateRequest struct {
		UUID     uuid.UUID `json:"-"`
		Name     string    `json:"name"`
		Login    string    `json:"login"`
		Password string    `json:"password"`
	}

	PatchRequest struct {
		UUID     uuid.UUID `json:"-"`
		Name     *string   `json:"name"`
		Login    *string   `json:"login"`
		Password *string   `json:"password"`
	}

	DeleteRequest struct {
		UUID uuid.UUID `json:"-"`
	}
)

type (
	GetAllResponse []GetResponse

	GetResponse struct {
		UUID  uuid.UUID `json:"uuid"`
		Name  string    `json:"name"`
		Login string    `json:"login"`
	}

	CreateResponse NoContent
	UpdateResponse NoContent
	PatchResponse  NoContent
	DeleteResponse NoContent
)
