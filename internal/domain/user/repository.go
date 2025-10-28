package user

type Lister interface {
	List(*ListRequest) (ListResponse, error)
}

type Getter interface {
	Get(*GetRequest) (Response, error)
}

type GetterByLogin interface {
	GetByLogin(*GetByLoginRequest) (Response, error)
}

type Creater interface {
	Create(*CreateRequest) (Response, error)
}

type Patcher interface {
	Patch(*PatchRequest) (Response, error)
}

type Deleter interface {
	Delete(*DeleteRequest) error
}

type Repository interface {
	Lister
	Getter
	GetterByLogin
	Creater
	Patcher
	Deleter
}
