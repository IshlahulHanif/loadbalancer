package hostpool

type Usecase struct {
	repo repository
}

type repository struct {
	hostpool RepoHostpool
}
