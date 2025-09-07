package user

import "github.com/alan-b-lima/prp/pkg/uuid"

type NoContent struct{}

type (
	// Has not been useful
	Request interface{ _request() }

	GetAllRequest NoContent

	GetRequest struct {
		UUID uuid.UUID `json:"uuid"`
	}

	CreateRequest struct {
		Name  string `json:"name"`
		Login string `json:"login"`
	}

	UpdateRequest struct {
		UUID  uuid.UUID `json:"uuid"`
		Name  string    `json:"name"`
		Login string    `json:"login"`
	}

	PatchRequest struct {
		UUID  uuid.UUID `json:"uuid"`
		Name  *string   `json:"name"`
		Login *string   `json:"login"`
	}

	DeleteRequest struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

func (*GetAllRequest) _request() {}
func (*GetRequest) _request()    {}
func (*CreateRequest) _request() {}
func (*UpdateRequest) _request() {}
func (*DeleteRequest) _request() {}

type (
	// Has not been useful
	Response interface{ _response() }

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

func (*GetAllResponse) _response() {}
func (*GetResponse) _response()    {}
func (*CreateResponse) _response() {}
func (*UpdateResponse) _response() {}
func (*DeleteResponse) _response() {}
