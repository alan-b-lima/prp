package session

func Get(repo Getter, req GetRequest) (Response, error) {
	return repo.Get(req.UUID)
}

func Create(repo Creater, req CreateRequest) (Response, error) {
	return repo.Create(req.User, req.MaxAge)
}
