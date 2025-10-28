package user

func List(repo Lister, req *ListRequest) (ListResponse, error) {
	res, err := repo.List(req)
	return res, err
}

func Get(repo Getter, req *GetRequest) (Response, error) {
	res, err := repo.Get(req)
	return res, err
}

func GetByLogin(repo GetterByLogin, req *GetByLoginRequest) (Response, error) {
	res, err := repo.GetByLogin(req)
	return res, err
}

func Create(repo Creater, req *CreateRequest) (Response, error) {
	res, err := repo.Create(req)
	return res, err
}

func Patch(repo Patcher, req *PatchRequest) (Response, error) {
	res, err := repo.Patch(req)
	return res, err
}

func Delete(repo Deleter, req *DeleteRequest) error {
	return repo.Delete(req)
}
